import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, switchMap, throwError} from 'rxjs';
import {RestResponse} from '../types/rest-response';
import {WebsocketService} from './websocket.service';

@Injectable({
  providedIn: 'root'
})
export class RestHelperService {
  http = inject(HttpClient)
  socketService = inject(WebsocketService)

  postAndAwaitProjectionUpdate<T, D>(url: string, data: D): Observable<T> {
    return this.http.post<RestResponse<string>>(url, data).pipe(
      switchMap(res => {
        if (res.errors.length > 0) {
          return throwError(() => new Error('Post-Fehler: ' + res.errors.join(', ')));
        }

        const aggregateId = res.data;
        const ws$ = this.socketService.listen(aggregateId);

        return new Observable<T>(observer => {
          const subscription = ws$.subscribe({
            next: message => {
              if (message.type === 'PROJECTION_UPDATED') {
                observer.next(res as unknown as T);
                observer.complete();
                subscription.unsubscribe();
              }
            },
            error: err => {
              observer.error(err);
            }
          });

          return () => subscription.unsubscribe();
        });
      })
    );
  }
}

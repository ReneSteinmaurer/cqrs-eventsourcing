import {inject, Injectable, signal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {MediumDetail} from '../types/medium-details';
import {Nutzer} from '../../../../shared/types/nutzer';
import {RestResponse} from '../../../../shared/types/rest-response';
import {WebsocketService} from '../../../../shared/services/websocket.service';
import {EMPTY, tap} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MediumDetailService {
  http = inject(HttpClient)
  socketService = inject(WebsocketService)
  #details = signal<MediumDetail | undefined>(undefined)
  details = this.#details.asReadonly()

  #nutzerPossibilities = signal<Nutzer[]>([])
  customerOptions = this.#nutzerPossibilities.asReadonly()

  loadDetails(mediumId: string) {
    this.http.get<MediumDetail>(`http://localhost:8080/bibliothek/medium?mediumId=${mediumId}`).subscribe(res => {
      this.#details.set(res);
    });
  }

  findNutzerByEmail(email: string) {
    this.http.get<RestResponse<Nutzer[]>>(`http://localhost:8080/nutzer/find-by-email?email=${email}`).subscribe(res => {
      this.#nutzerPossibilities.set(res.data);
    })
  }

  verleihen(nutzerId: string) {
    const mediumId = this.#details()?.mediumId
    if (!mediumId) {
      console.error('medium id is nullish!')
      return EMPTY
    }
    return this.http.post<RestResponse<string>>('http://localhost:8080/bibliothek/verleihe-medium', {nutzerId, mediumId}).pipe(tap(res => {
      if (res.errors.length > 0) {
        // todo handle errors
        return
      }
      this.socketService.listen(res.data).subscribe(() => {
        this.loadDetails(mediumId)
      });
    }));
  }

  zurueckgeben() {
    const mediumId = this.#details()?.mediumId
    const nutzerId = this.#details()?.verliehenNutzerId
    if (!mediumId || !nutzerId) {
      console.error('mediumId or nutzerId is nullish!')
      return EMPTY
    }
    return this.http.post<RestResponse<string>>('http://localhost:8080/bibliothek/gebe-medium-zurueck', {nutzerId, mediumId}).pipe(tap(res => {
      if (res.errors.length > 0) {
        // todo handle errors
        return
      }
      this.socketService.listen(res.data).subscribe(() => {
        this.loadDetails(mediumId)
      });
    }));
  }
}

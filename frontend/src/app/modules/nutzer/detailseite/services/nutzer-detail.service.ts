import {inject, Injectable, signal} from '@angular/core';
import {NutzerDetails} from '../types/nutzer-details';
import {HttpClient} from '@angular/common/http';
import {RestHelperService} from '../../../../shared/services/rest-helper.service';

@Injectable({
  providedIn: 'root'
})
export class NutzerDetailService {
  http = inject(HttpClient)
  rest = inject(RestHelperService)
  #details = signal<NutzerDetails | undefined>(undefined)
  details = this.#details.asReadonly()

  loadDetails(nutzerId: string) {
    this.http.get<NutzerDetails>(`http://localhost:8080/nutzer?nutzerId=${nutzerId}`).subscribe(res => {
      this.#details.set(res);
    });
  }
}

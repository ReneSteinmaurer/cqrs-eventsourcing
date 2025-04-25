import {inject, Injectable, signal} from '@angular/core';
import {MediumBestand} from '../../shared/types/medium-bestand';
import {HttpClient} from '@angular/common/http';

type WebSocketSubscribeMessage = {
  type: 'subscribe'
  topic: string
}

@Injectable({
  providedIn: 'root'
})
export class LandingPageService {
  private http = inject(HttpClient)
  medien = signal<MediumBestand[]>([])

  getAllMedien() {
    this.http.get<MediumBestand[]>('http://localhost:8080/bibliothek/bestand').subscribe(res => {
      this.medien.set(res)
    });
  }

  constructor() {
    console.log('init')
  }
}

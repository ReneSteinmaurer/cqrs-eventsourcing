import {Injectable} from '@angular/core';
import {httpResource} from '@angular/common/http';
import {MediumBestand} from '../../shared/types/medium-bestand';

@Injectable({
  providedIn: 'root'
})
export class LandingPageService {
  medien = httpResource<MediumBestand[]>(() => ({
    url: 'http://localhost:8080/bibliothek/bestand',
    method: 'GET',
  }));
}

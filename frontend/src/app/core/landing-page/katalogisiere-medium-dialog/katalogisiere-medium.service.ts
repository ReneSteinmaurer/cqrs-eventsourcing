import {inject, Injectable, signal} from '@angular/core';
import {HttpClient, httpResource} from '@angular/common/http';
import {RestResponse} from '../../../shared/types/rest-response';
import {LandingPageService} from '../landing-page.service';
import {WebSocketMessage, WebsocketService} from '../../../shared/services/websocket.service';
import {RestHelperService} from '../../../shared/services/rest-helper.service';
import {tap} from 'rxjs';

export type KatalogisiereCommand = {
  MediumId: string
  ISBN: string
  Signature: string
  Standort: string
}

@Injectable({
  providedIn: 'root'
})
export class KatalogisiereMediumService {
  private restHelper = inject(RestHelperService)
  private landingPageService = inject(LandingPageService)

  katalogisiereMedium(cmd: KatalogisiereCommand) {
    this.restHelper.postAndAwaitProjectionUpdate('http://localhost:8080/bibliothek/katalogisiere-medium', cmd).subscribe({
      next: () => {
        this.landingPageService.medien.reload();
      },
      error: (err) => {
        console.error('Fehler beim Verleihen:', err);
      }
    });
  }
}

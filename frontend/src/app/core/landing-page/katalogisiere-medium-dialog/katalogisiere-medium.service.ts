import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {RestResponse} from '../../../shared/types/rest-response';
import {LandingPageService} from '../landing-page.service';
import {WebSocketMessage, WebsocketService} from '../../../shared/services/websocket.service';

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
  private http = inject(HttpClient)
  private websocketService = inject(WebsocketService)
  landingPageService = inject(LandingPageService)

  katalogisiereMedium(cmd: KatalogisiereCommand) {
    this.http.post<RestResponse<string>>('http://localhost:8080/bibliothek/katalogisiere-medium', {...cmd}).subscribe(res => {
      this.websocketService.listen(res.data).subscribe((msg) => {
          if (msg.type === 'PROJECTION_UPDATED') {
            this.landingPageService.getAllMedien()
          }
      });
    });
  }
}

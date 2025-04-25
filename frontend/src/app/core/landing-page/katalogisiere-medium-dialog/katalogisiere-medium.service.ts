import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {RestResponse} from '../../../shared/types/rest-response';
import {LandingPageService} from '../landing-page.service';

export type KatalogisiereCommand = {
  MediumId: string
  ISBN: string
  Signature: string
  Standort: string
}

const WebSocketMessageTypeKeys = [
  'PROJECTION_UPDATED'
] as const

type WebSocketMessageType = typeof WebSocketMessageTypeKeys[number]

export type WebSocketMessage<T> = {
  type: WebSocketMessageType
  data: T
}

@Injectable({
  providedIn: 'root'
})
export class KatalogisiereMediumService {
  private http = inject(HttpClient)
  private socket?: WebSocket
  landingPageService = inject(LandingPageService)

  katalogisiereMedium(cmd: KatalogisiereCommand) {
    this.http.post<RestResponse<string>>('http://localhost:8080/bibliothek/katalogisiere-medium', {...cmd}).subscribe(res => {
      console.log(res.data)
      this.openWebsocketConnection(res.data)
    });
  }

  openWebsocketConnection(aggregateId: string) {
    console.log('INIT WEBSOCKET CONNECTION')
    this.socket = new WebSocket(`ws://localhost:8080/ws?aggregateId=${aggregateId}`)

    this.socket.onmessage = (event) => {
      const message = JSON.parse(event.data) as WebSocketMessage<string>
      console.log(message)
      if (message.type === 'PROJECTION_UPDATED') {
        this.landingPageService.getAllMedien()
      }
    }
  }
}

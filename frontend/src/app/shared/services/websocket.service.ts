import { Injectable } from '@angular/core';
import {webSocket, WebSocketSubject} from 'rxjs/webSocket';
import {Observable} from 'rxjs';

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
export class WebsocketService {
  private wsSubject: WebSocketSubject<WebSocketMessage<string>> | null = null;
  private socket?: WebSocket

  listen(aggregateId: string): Observable<WebSocketMessage<string>>{
    this.wsSubject = webSocket<WebSocketMessage<string>>({
      url: `ws://localhost:8080/ws?aggregateId=${aggregateId}`,
      deserializer: msg => JSON.parse(msg.data)
    })
    return this.wsSubject.asObservable();
  }

}

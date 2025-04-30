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
  private connections = new Map<string, WebSocketSubject<WebSocketMessage<string>>>();

  listen(aggregateId: string): Observable<WebSocketMessage<string>> {
    const existing = this.connections.get(aggregateId);

    if (existing && !existing.closed) {
      return existing;
    }

    const ws = webSocket<WebSocketMessage<string>>({
      url: `ws://localhost:8080/ws?aggregateId=${aggregateId}`,
      deserializer: msg => JSON.parse(msg.data)
    });

    this.connections.set(aggregateId, ws);
    return ws;
  }

  close(aggregateId: string): void {
    const ws = this.connections.get(aggregateId);
    if (ws && !ws.closed) {
      ws.complete();
    }
    this.connections.delete(aggregateId);
  }

}

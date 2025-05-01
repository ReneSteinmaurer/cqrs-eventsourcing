import {Component, computed, inject, input} from '@angular/core';
import {RowLabelComponent} from '../row-label/row-label.component';
import {EventType, eventTypeToIcon, HistoryEvent, isEventType} from '../../types/history-events';
import {DatePipe} from '@angular/common';
import {Router} from '@angular/router';

@Component({
  selector: 'app-event-history',
  imports: [
    RowLabelComponent,
    DatePipe
  ],
  templateUrl: './event-history.component.html',
})
export class EventHistoryComponent {
  router = inject(Router)
  eventTypes = input<EventType[]>([])
  historyEvents = input.required<HistoryEvent[]>()
  events = computed(() => {
    if (this.eventTypes().length === 0) {
      return this.historyEvents();
    }
    return this.historyEvents().filter(val => this.eventTypes().includes(val.eventType));
  });

  openNutzerDetailseite(id: string) {
    this.router.navigate(['/nutzer', id]);
  }

  openMediumDetailseite(id: string) {
    this.router.navigate(['/medium', id]);
  }

  getIcon(event: HistoryEvent): string {
    return eventTypeToIcon[event.eventType as keyof typeof eventTypeToIcon] ?? 'help';
  }

  protected readonly isEventType = isEventType;
}

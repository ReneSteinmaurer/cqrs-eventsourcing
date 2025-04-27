import {Component, input} from '@angular/core';
import {RowLabelComponent} from '../../../../../shared/ui/row-label/row-label.component';
import {
  HistoryEvent,
  MediumErworbenEvent,
  MediumKatalogisiertEvent,
  MediumVerliehenEvent,
  MediumZurueckgegebenEvent
} from '../../types/medium-details';
import {DatePipe} from '@angular/common';

@Component({
  selector: 'app-history',
  imports: [
    RowLabelComponent,
    DatePipe
  ],
  template: `
    <div class="flex flex-col gap-2">
      @for (event of historyEvents(); track $index) {
        <div class="bg-[#00dddd]/10 rounded-xl p-4 shadow-md hover:bg-[#00dddd]/15 transition">
          <app-row-label
            [icon]="
          isErworben(event) ? 'library_add' :
          isKatalogisiert(event) ? 'label_important' :
          'assignment_returned'
        "
            [iconClass]="'text-primary'"
            [label]="(event.timestamp | date:'dd.MM.yyyy') ?? ''"
          >
            <div class="flex flex-col">

              @if (isErworben(event)) {
                <span class="font-semibold text-gray-100 text-sm">
              Medium erworben: {{ event.payload.Name }} ({{ event.payload.MediumType }})
            </span>
              }

              @if (isKatalogisiert(event)) {
                <span class="font-semibold text-gray-100 text-sm">
              Katalogisiert: {{ event.payload.Standort }}, Signatur: {{ event.payload.Signature }}
            </span>
              }

              @if (isVerliehen(event)) {
                <span class="font-semibold text-gray-100 text-sm">
              Verliehen an <span class="hover:underline hover:cursor-pointer">{{ event.payload.NutzerId }}</span>
            </span>
                @if (event.payload.Von && event.payload.Bis) {
                  <span class="text-gray-400 text-xs mt-1">
                von {{ event.payload.Von | date:'dd.MM.yyyy' }} bis {{ event.payload.Bis | date:'dd.MM.yyyy' }}
              </span>
                }
              }

              @if (isZurueckgegeben(event)) {
                <span class="font-semibold text-gray-100 text-sm">
              Medium zurückgegeben von <span class="hover:underline hover:cursor-pointer">{{ event.payload.NutzerId }}</span>
            </span>
                @if (event.payload.Date) {
                  <span class="text-gray-400 text-xs mt-1">
                am {{ event.payload.Date | date:'dd.MM.yyyy' }}
              </span>
                }
              }

            </div>
          </app-row-label>
        </div>
      }
    </div>
  `,
})
export class HistoryComponent {
  historyEvents = input.required<HistoryEvent[]>()

  isErworben(event: HistoryEvent): event is MediumErworbenEvent {
    return event.eventType === 'MediumErworbenEvent';
  }

  isKatalogisiert(event: HistoryEvent): event is MediumKatalogisiertEvent {
    return event.eventType === 'MediumKatalogisiertEvent';
  }

  isVerliehen(event: HistoryEvent): event is MediumVerliehenEvent {
    return event.eventType === 'MediumVerliehenEvent';
  }

  isZurueckgegeben(event: HistoryEvent): event is MediumZurueckgegebenEvent {
    return event.eventType === 'MediumZurueckgegebenEvent';
  }

}

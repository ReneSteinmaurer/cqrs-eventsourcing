import {Component, computed, inject} from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';
import {ActivatedRoute, Router} from '@angular/router';
import {CardComponent} from '../../../shared/ui/card/card.component';
import {RowLabelComponent} from '../../../shared/ui/row-label/row-label.component';
import {AusgeliehenComponent} from './components/ausgeliehen/ausgeliehen.component';
import {NutzerDetailService} from './services/nutzer-detail.service';
import {VerlorenComponent} from './components/verloren/verloren.component';
import {DatePipe, Location} from '@angular/common';
import {EventHistoryComponent} from '../../../shared/ui/event-history/event-history.component';
import {EventType} from '../../../shared/types/history-events';
import {MatTabGroup, MatTabsModule} from '@angular/material/tabs';

@Component({
  selector: 'app-detailseite',
  imports: [
    MatButton,
    MatIcon,
    CardComponent,
    RowLabelComponent,
    AusgeliehenComponent,
    VerlorenComponent,
    DatePipe,
    EventHistoryComponent,
    MatTabGroup,
    MatTabsModule
  ],
  template: `
    <div class="flex justify-center items-center">
      <div class="p-6 space-y-6 w-full lg:w-5/6 xl:w-1/2 mt-12">
        <div class="flex justify-between">
          <button (click)="navigateToPreviousPage()" mat-stroked-button color="primary" class="mb-4">
            <mat-icon>arrow_back</mat-icon>
            Zurück
          </button>
        </div>
        <!--<div class="border border-solid border-red-400 bg-red-400/30 p-4 rounded-lg">
          <div class="flex items-center gap-4">
            <mat-icon class="text-red-400">warning</mat-icon>
            Gesperrt bis 15.05.2025 wegen unbezahlter Mahngebühren.
          </div>
        </div>-->
        <div class="grid gap-4">
          <app-card [title]="nutzerName()">
            <app-row-label label="Email">{{ details()?.nutzerId }}</app-row-label>
            <app-row-label label="Mitglied seit">{{ details()?.registriertAm | date: 'dd.MM.YYYY' }}</app-row-label>
          </app-card>
          <mat-tab-group mat-stretch-tabs [dynamicHeight]="false" class="min-h-[400px]">
            @if (verloreneMedien().length > 0) {
              <mat-tab label="Verlorene">
                <div class="flex p-4 flex-col gap-2 max-h-72 overflow-y-auto pr-2">
                  @for (verloren of verloreneMedien(); track $index) {
                    <app-verloren [verloren]="verloren"/>
                  }
                </div>
              </mat-tab>
            }
            <mat-tab label="Ausgeliehen">
              <div class="flex flex-col p-4 gap-2 max-h-72 overflow-y-auto pr-2">
                @for (ausgeliehen of details()?.aktiveAusleihen ?? []; track $index) {
                  <app-ausgeliehen [ausgeliehen]="ausgeliehen"/>
                } @empty {
                  <p>Noch kein Medium ausgeliehen</p>
                }
              </div>
            </mat-tab>
            <mat-tab label="Historie">
              <div class="flex flex-col p-4 gap-2 max-h-72 overflow-y-auto pr-2">
                <app-event-history [eventTypes]="eventTypes" [historyEvents]="history()"/>
              </div>
            </mat-tab>
          </mat-tab-group>
        </div>
      </div>
    </div>`,
  styles: ``
})
export class DetailseiteComponent {
  router = inject(Router)
  route = inject(ActivatedRoute)
  location = inject(Location);
  detailService = inject(NutzerDetailService)
  details = this.detailService.details
  history = computed(() => {
    return this.details()?.historie ?? [];
  });
  verloreneMedien = computed(() => {
    return this.details()?.verloreneMedien ?? [];
  })
  nutzerName = computed(() => {
    return this.details()?.vorname + ' ' + this.details()?.nachname;
  });

  eventTypes: EventType[] = [
    'NutzerRegistriertEvent',
    'MediumVerliehenEvent',
    "MediumVerlorenDurchBenutzerEvent",
    "MediumZurueckgegebenEvent",
    "MediumWiederaufgefundenDurchNutzerEvent"
  ]

  constructor() {
    this.route.params.subscribe(params => {
      const nutzerId = params['id']
      if (!nutzerId) {
        return
      }
      this.detailService.loadDetails(nutzerId)
    })
  }

  navigateToPreviousPage() {
    this.location.back()
  }

}

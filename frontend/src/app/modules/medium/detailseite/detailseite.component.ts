import {Component, inject, input} from "@angular/core";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatIconModule} from "@angular/material/icon";
import {MatTabsModule} from "@angular/material/tabs";
import {MatBadgeModule} from '@angular/material/badge';
import {RowLabelComponent} from '../../../shared/ui/row-label/row-label.component';
import {CardComponent} from '../../../shared/ui/card/card.component';
import {ActivatedRoute, Router} from '@angular/router';
import {StatusChipComponent} from './components/status-chip/status-chip.component';
import {MediumDetailService} from './services/medium-detail.service';
import {EventHistoryComponent} from '../../../shared/ui/event-history/event-history.component';
import {DatePipe, Location} from '@angular/common';
import {AktionenComponent} from './components/aktionen/aktionen.component';
import {MatDialog} from '@angular/material/dialog';
import {VerleihenComponent} from './components/aktionen/verleihen/verleihen.component';
import {ZuruecknehmenDialogComponent} from './components/aktionen/zuruecknehmen/zuruecknehmen-dialog.component';
import {VerlorenDurchNutzerComponent} from './components/aktionen/verloren-durch-nutzer/verloren-durch-nutzer-dialog.component';
import {BestandsverlustDialogComponent} from './components/aktionen/bestandsverlust-dialog/bestandsverlust-dialog.component';
import {
  BestandsverlustAufhebenDialogComponent
} from './components/aktionen/bestandsverlust-aufheben-dialog/bestandsverlust-aufheben-dialog.component';
import {
  WiederaufgefundenDurchNutzerDialogComponent
} from './components/aktionen/wiederaufgefunden-durch-nutzer-dialog/wiederaufgefunden-durch-nutzer-dialog.component';

@Component({
  selector: 'app-detailseite',
  imports: [
    MatCardModule,
    MatButtonModule,
    MatTabsModule,
    MatIconModule,
    MatBadgeModule,
    RowLabelComponent,
    CardComponent,
    StatusChipComponent,
    EventHistoryComponent,
    DatePipe,
    AktionenComponent,
  ],
  template: `
    @if (mediumDetails()) {
      <div class="flex justify-center items-center">
        <div class="p-6 space-y-6 w-full lg:w-5/6 xl:w-1/2 mt-12">
          <div class="flex justify-between">
            <button (click)="navigateToPreviousPage()" mat-stroked-button color="primary" class="mb-4">
              <mat-icon>arrow_back</mat-icon>
              Zurück
            </button>
            <app-aktionen (verlorenDurchNutzer)="openVerlorenDurchNutzerDialog()"
                          (zuruecknehmen)="openZuruecknehmenDialog()"
                          (verleihen)="openVerleihenDialog()"
                          (wiederaufgefunden)="openWiederaufgefundenDialog()"
                          (bestandsverlust)="openBestandsverlustDialog()"
                          [aktuellVerliehen]="mediumDetails()?.aktuellVerliehen ?? false"
                          [verloren]="mediumDetails()?.status === 'VERLOREN'"/> <!-- do better -->
          </div>

          <div class="space-x-6">
            <app-card [title]="mediumDetails()!.titel">
              <div class="absolute top-4 right-4">
                <app-status-chip [status]="mediumDetails()!.status"></app-status-chip>
              </div>
              <app-row-label label="ISBN">{{ mediumDetails()?.isbn }}</app-row-label>
              <app-row-label label="Genre">{{ mediumDetails()?.genre }}</app-row-label>
              <app-row-label label="Typ">{{ mediumDetails()?.typ }}</app-row-label>
            </app-card>

            <app-card [title]="'📍 Standortinformationen'">
              <app-row-label label="Standort">{{ mediumDetails()?.standort }}</app-row-label>
              <app-row-label label="Signatur">{{ mediumDetails()?.signatur }}</app-row-label>
              <app-row-label label="Exemplar-Code">{{ mediumDetails()?.exemplarCode }}
              </app-row-label>
            </app-card>

            @if (mediumDetails()?.aktuellVerliehen) {
              <app-card [title]="'📦 Verleihstatus'">
                <app-row-label label="Aktuell verliehen an">{{ mediumDetails()?.verliehenAn }}</app-row-label>
                <app-row-label label="Verliehen von">{{ mediumDetails()?.verliehenVon | date: 'dd-MM-YYYY' }}
                </app-row-label>
                <app-row-label label="Fällig bis">{{ mediumDetails()?.faelligBis | date: 'dd-MM-YYYY' }}</app-row-label>
              </app-card>
            }
          </div>

          <mat-tab-group mat-stretch-tabs [dynamicHeight]="false" class="min-h-[400px]">
            <mat-tab label="Details">
              <div class="p-4 space-y-4">
                <app-card [title]="'📖 Medium-Informationen'">
                  <app-row-label label="Erworben am">{{ mediumDetails()?.erworbenAm | date: 'dd-MM-YYYY' }}
                  </app-row-label>
                  <app-row-label label="Katalogisiert am">{{ mediumDetails()?.katalogisiertAm | date: 'dd-MM-YYYY' }}
                  </app-row-label>
                  <app-row-label label="Aktueller Standort">{{ mediumDetails()?.standort }}</app-row-label>
                  <app-row-label label="Exemplar-Code">{{ mediumDetails()?.exemplarCode }}</app-row-label>
                </app-card>
              </div>
            </mat-tab>

            <mat-tab label="Historie">
              <div class="p-4 space-y-4">
                <div class="flex flex-col space-y-4 max-h-[300px] overflow-y-auto pr-2">
                  <app-event-history [historyEvents]="mediumDetails()?.historie ?? []"/>
                </div>
              </div>
            </mat-tab>
          </mat-tab-group>

          <div class="flex flex-wrap justify-end gap-4 mt-4">

          </div>
        </div>
      </div>
    }
  `,
  styles: ``
})
export class DetailseiteComponent {
  router = inject(Router)
  route = inject(ActivatedRoute)
  location = inject(Location)
  dialog = inject(MatDialog)
  detailService = inject(MediumDetailService)
  mediumId = input.required<string>({alias: 'mediumId'})
  mediumDetails = this.detailService.details

  constructor() {
    this.route.params.subscribe(params => {
      const mediumId = params['id']
      if (!mediumId) {
        return
      }
      this.detailService.loadDetails(mediumId)
    })
  }

  navigateToPreviousPage() {
    this.location.back()
  }

  openVerleihenDialog() {
    this.dialog.open(VerleihenComponent, {
      width: '70%',
      height: '40%',
    })
  }

  openZuruecknehmenDialog() {
    this.dialog.open(ZuruecknehmenDialogComponent, {
      width: '70%',
      height: '20%',
    })
  }

  openVerlorenDurchNutzerDialog() {
    this.dialog.open(VerlorenDurchNutzerComponent, {
      width: '70%',
      height: '20%',
    })
  }

  openBestandsverlustDialog() {
    this.dialog.open(BestandsverlustDialogComponent, {
      width: '70%',
      height: '20%',
    })
  }

  openWiederaufgefundenDialog() {
    if (this.mediumDetails()?.verlorenVonNutzerId) {
      console.log('wiederaufgefunden')
      this.dialog.open(WiederaufgefundenDurchNutzerDialogComponent, {
        width: '70%',
        height: '20%',
      })
      return
    }
    this.dialog.open(BestandsverlustAufhebenDialogComponent, {
      width: '70%',
      height: '20%',
    })
  }
}

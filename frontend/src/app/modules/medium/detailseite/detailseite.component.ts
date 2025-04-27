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
import {HistoryComponent} from './components/history/history.component';
import {DatePipe} from '@angular/common';
import {AktionenComponent} from './components/aktionen/aktionen.component';
import {MatDialog} from '@angular/material/dialog';
import {VerleihenComponent} from './components/verleihen/verleihen.component';
import {ZuruecknehmenComponent} from './components/zuruecknehmen/zuruecknehmen.component';

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
    HistoryComponent,
    DatePipe,
    AktionenComponent,
  ],
  template: `
    @if (mediumDetails()) {
      <div class="flex justify-center items-center">
        <div class="p-6 space-y-6 w-full lg:w-5/6 xl:w-1/2 mt-12">
          <div class="flex justify-between">
            <button (click)="navigateToLandingPage()" mat-stroked-button color="primary" class="mb-4">
              <mat-icon>arrow_back</mat-icon>
              Zurück zur Übersicht
            </button>
            <app-aktionen (zuruecknehmen)="openZuruecknehmenDialog()" (verleihen)="openVerleihenDialog()"
                          [aktuellVerliehen]="mediumDetails()?.aktuellVerliehen ?? false"/>
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
                  <app-history [historyEvents]="mediumDetails()?.historie ?? []"/>
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

  navigateToLandingPage() {
    this.router.navigate([''])
  }

  editStandort() {

  }

  editBuchdaten() {

  }

  openVerleihenDialog() {
    this.dialog.open(VerleihenComponent, {
      width: '70%',
      height: '40%',
    })
  }

  openZuruecknehmenDialog() {
    this.dialog.open(ZuruecknehmenComponent, {
      width: '70%',
      height: '20%',
    })
  }
}

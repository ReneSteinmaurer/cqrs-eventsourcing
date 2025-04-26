import {Component, inject} from "@angular/core";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatIconModule} from "@angular/material/icon";
import {MatTabsModule} from "@angular/material/tabs";
import {MatBadgeModule} from '@angular/material/badge';
import {RowLabelComponent} from '../../../shared/ui/row-label/row-label.component';
import {CardComponent} from '../../../shared/ui/card/card.component';
import {Router} from '@angular/router';
import {StatusChipComponent} from './status-chip/status-chip.component';

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
  ],
  template: `
    <div class="flex justify-center items-center">
      <div class="p-6 space-y-6 w-full lg:w-5/6 xl:w-1/2">

        <button (click)="navigateToLandingPage()" mat-stroked-button color="primary" class="mt-12 mb-4">
          <mat-icon>arrow_back</mat-icon>
          Zurück zur Übersicht
        </button>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="col-span-1 md:col-span-2">
            <app-card [title]="'Sapiens: A Brief History of Humankind'">
              <div class="absolute top-4 right-4">
                <app-status-chip [status]="'verliehen'"></app-status-chip>
              </div>
              <div class="space-y-2">
                <p><strong>ISBN:</strong> 978-0-14-312779-9</p>
                <p><strong>Genre:</strong> Sachbuch</p>
                <p><strong>Typ:</strong> Buch</p>
              </div>
            </app-card>
          </div>

          <app-card [title]="'📍 Standortinformationen'">
            <div class="space-y-2">
              <p><strong>Standort:</strong> Hauptbibliothek Wien</p>
              <p><strong>Signatur:</strong> FIC-REI-2025</p>
              <p><strong>Exemplar-Code:</strong> EX123456789</p>
            </div>
          </app-card>

          <app-card [title]="'📦 Verleihstatus'">
            <div class="space-y-2">
              <p><strong>Aktuell verliehen an:</strong> -</p>
              <p><strong>Verliehen von:</strong> -</p>
              <p><strong>Fällig bis:</strong> -</p>
            </div>
          </app-card>
        </div>

        <mat-tab-group mat-stretch-tabs [dynamicHeight]="false" class="min-h-[400px]">
          <mat-tab label="Details">
            <div class="p-4 space-y-4">
              <app-card [title]="'📖 Medium-Informationen'">
                <div class="space-y-2">
                  <p><strong>Erworben am:</strong> 21.04.2025</p>
                  <p><strong>Katalogisiert am:</strong> 21.04.2025</p>
                  <p><strong>Aktueller Standort:</strong> Hauptbibliothek Wien</p>
                  <p><strong>Exemplar-Code:</strong> EX123456789</p>
                </div>
              </app-card>
            </div>
          </mat-tab>

          <mat-tab label="Historie">
            <div class="p-4 space-y-4">
              <div class="flex flex-col space-y-4 max-h-[300px] overflow-y-auto pr-2">
                <app-row-label
                  [icon]="'library_add'"
                  [iconClass]="'text-primary'"
                  [label]="'21.04.2025'"
                >
                  Medium erworben
                </app-row-label>

                <app-row-label
                  [icon]="'label_important'"
                  [iconClass]="'text-primary'"
                  [label]="'21.04.2025'"
                >
                  Medium katalogisiert
                </app-row-label>

                <app-row-label
                  [icon]="'assignment_returned'"
                  [iconClass]="'text-accent'"
                  [label]="'22.04.2025'"
                >
                  Medium verliehen an Rene
                </app-row-label>

                <app-row-label
                  [icon]="'assignment_turned_in'"
                  [iconClass]="'text-warn'"
                  [label]="'22.04.2025'"
                >
                  Medium zurückgegeben
                </app-row-label>
              </div>
            </div>
          </mat-tab>
        </mat-tab-group>

        <div class="flex flex-wrap justify-end gap-4 mt-4">
          <button mat-flat-button color="primary" class="px-6">
            Medium verleihen
          </button>
          <button mat-flat-button color="warn" class="px-6">
            Medium zurücknehmen
          </button>
          <button mat-stroked-button color="accent" class="px-6">
            Bearbeiten
          </button>
        </div>

      </div>
    </div>
  `,
  styles: ``
})
export class DetailseiteComponent {
  router = inject(Router)

  navigateToLandingPage() {
    this.router.navigate([''])
  }
}

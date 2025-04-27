import {Component, input, output} from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';
import {MatIcon} from '@angular/material/icon';

@Component({
  selector: 'app-aktionen',
  imports: [
    MatButton,
    MatIcon,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger
  ],
  template: `
    <div class="flex gap-4">
      @if (!aktuellVerliehen()) {
        <button (click)="verleihen.emit()" mat-flat-button class="px-6">
          Medium verleihen
        </button>
      }
      @if (aktuellVerliehen()) {
        <button (click)="zuruecknehmen.emit()" mat-flat-button class="px-6">
          Medium zurückgegeben
        </button>
      }
      <button mat-stroked-button [matMenuTriggerFor]="editMenu" class="px-6">
        <mat-icon>edit</mat-icon>
        Bearbeiten
      </button>

      <mat-menu #editMenu="matMenu">
        <button mat-menu-item (click)="editStandort.emit()">
          <mat-icon>location_on</mat-icon>
          <span>Standort gewechselt</span>
        </button>

        <button mat-menu-item (click)="editBuchdaten.emit()">
          <mat-icon>menu_book</mat-icon>
          <span>Buchdaten korrigieren</span>
        </button>
      </mat-menu>
    </div>
  `,
  styles: ``
})
export class AktionenComponent {
  aktuellVerliehen = input.required<boolean>()
  editStandort = output<void>()
  editBuchdaten = output<void>()
  verleihen = output<void>()
  zuruecknehmen = output<void>()

}

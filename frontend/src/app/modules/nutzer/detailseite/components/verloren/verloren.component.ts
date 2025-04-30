import {Component, inject, input} from '@angular/core';
import {DatePipe} from '@angular/common';
import {VerlorenesMedium} from '../../types/nutzer-details';
import {MatIcon} from '@angular/material/icon';
import {Router} from '@angular/router';

@Component({
  selector: 'app-verloren',
  imports: [
    DatePipe,
    MatIcon
  ],
  template: `
    <div class="p-4 rounded-xl bg-red-400/10 hover:bg-red-400/15 transition shadow-sm">
      <div class="flex items-center justify-between">
        <div>
          <h3><span (click)="openMediumDetailseite()"
                    class="hover:underline hover:cursor-pointer">{{ verloren().titel }}</span></h3>
          <p>
            Typ: <span class="font-medium text-white">TODO (BUCH)</span>
          </p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center md:justify-between text-sm gap-3">
        <div class="text-gray-300">
          <span class="font-medium text-white">Ausgeliehen am:</span>
          {{ verloren().ausgeliehenAm | date: 'dd.MM.yyyy HH:mm' }}
        </div>
        <div>
           <span
             class="inline-flex items-center gap-1 bg-red-600/20 text-red-400 px-3 py-1 rounded-full text-xs font-semibold">
             <mat-icon class="text-red-400 text-sm">error</mat-icon>
             Verloren
           </span>
        </div>
      </div>
    </div>
  `,
})
export class VerlorenComponent {
  router = inject(Router)
  verloren = input.required<VerlorenesMedium>()

  openMediumDetailseite() {
    this.router.navigate(['medium', this.verloren().mediumId])
  }
}

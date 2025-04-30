import {Component, inject, input} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {AktiveAusleihe, VerlorenesMedium} from '../../types/nutzer-details';
import {DatePipe} from '@angular/common';
import {Router} from '@angular/router';

@Component({
  selector: 'app-ausgeliehen',
  imports: [
    MatIcon,
    DatePipe
  ],
  template: `
    <div class="p-4 rounded-xl bg-[#00dddd]/10 hover:bg-[#00dddd]/15 transition shadow-sm">
      <div class="flex items-center justify-between">
        <div>
          <h3><span (click)="openMediumDetailseite()"
                    class="hover:underline hover:cursor-pointer">{{ ausgeliehen().titel }}</span></h3>
          <p>
            Typ: <span class="font-medium text-white">TODO (BUCH)</span>
          </p>
        </div>
        <mat-icon class="text-primary text-xl">history</mat-icon>
      </div>

      <div class="flex flex-col md:flex-row md:items-center md:justify-between text-sm gap-3">
        <div class="text-gray-300">
          <span class="font-medium text-white">Ausgeliehen am:</span>
          {{ ausgeliehen().ausgeliehenAm | date: 'dd.MM.yyyy HH:mm' }}
        </div>
        <div>
           <span
             class="inline-flex items-center gap-1 bg-green-600/20 text-green-400 px-3 py-1 rounded-full text-xs font-semibold">
             <mat-icon class="text-green-400 text-sm">check_circle</mat-icon>
             Aktiv
           </span>
        </div>
      </div>
    </div>
  `,
})
export class AusgeliehenComponent {
  router = inject(Router)
  ausgeliehen = input.required<AktiveAusleihe>()

  openMediumDetailseite() {
    this.router.navigate(['medium', this.ausgeliehen().mediumId])
  }
}

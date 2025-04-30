import {Component, inject} from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';
import {Router} from '@angular/router';

@Component({
  selector: 'app-detailseite',
  imports: [
    MatButton,
    MatIcon
  ],
  template: `
    <div class="flex justify-center items-center">
      <div class="p-6 space-y-6 w-full lg:w-5/6 xl:w-1/2 mt-12">
        <div class="flex justify-between">
          <button (click)="navigateToLandingPage()" mat-stroked-button color="primary" class="mb-4">
            <mat-icon>arrow_back</mat-icon>
            Zurück zur Übersicht
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

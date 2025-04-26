import {Component, input} from '@angular/core';
import {MatCard, MatCardContent, MatCardTitle} from '@angular/material/card';

@Component({
  selector: 'app-card',
  imports: [
    MatCardTitle,
    MatCard,
    MatCardContent
  ],
  template: `
    <mat-card class="p-6 space-y-4 shadow-md">
      @if (title()) {
        <mat-card-title class="text-2xl pb-4 font-semibold">{{ title() }}</mat-card-title>
      }
      <mat-card-content class="space-y-4">
        <ng-content></ng-content>
      </mat-card-content>
    </mat-card>
  `,
  styles: ``
})
export class CardComponent {
  title = input<string>('');
}

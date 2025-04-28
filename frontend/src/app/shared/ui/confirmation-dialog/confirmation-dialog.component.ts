import {Component, contentChild, effect, ElementRef, input, output} from '@angular/core';
import {MatButton} from '@angular/material/button';

@Component({
  selector: 'confirmation-dialog',
  imports: [
    MatButton
  ],
  template: `
    <div class="p-6 max-w-full max-h-full overflow-hidden">
      <h3 class="text-xl font-semibold m-4">{{ titel() }}</h3>
      <p class="text-sm text-gray-400 mt-4 text-center">
        <ng-content #subtextElem></ng-content>
      </p>
      @if (!subtext()) {
        <p class="text-sm text-gray-400 mt-4 text-center">
          Möchtest du wirklich fortfahren?
        </p>
      }
      <div class="flex justify-end gap-4 w-full pt-4">
        <button mat-button (click)="cancel.emit()">Abbrechen</button>
        <button mat-flat-button color="warn" (click)="confirm.emit()">Bestätigen</button>
      </div>
    </div>
  `,
  styles: ``
})
export class ConfirmationDialogComponent {
  subtext = contentChild<Element>('content')
  titel = input.required<string>()
  cancel = output<void>()
  confirm = output<void>()

}

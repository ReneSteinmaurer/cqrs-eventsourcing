import { Component, input } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import {NgClass} from '@angular/common';

@Component({
  selector: 'app-row-label',
  imports: [
    MatIcon,
    NgClass
  ],
  template: `
    <div class="flex items-start gap-4">
      @if (icon()) {
        <mat-icon [ngClass]="iconClass()">{{ icon() }}</mat-icon>
      }
      <div>
        <p><strong>{{ label() }}:</strong>
          <ng-content></ng-content>
        </p>
      </div>
    </div>
  `,
  styles: ``
})
export class RowLabelComponent {
  icon = input('')
  iconClass = input<string>('text-primary');
  label = input.required<string>()

}

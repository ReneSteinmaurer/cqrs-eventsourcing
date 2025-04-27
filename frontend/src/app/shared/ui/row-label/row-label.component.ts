import { Component, input} from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import {NgClass} from '@angular/common';

@Component({
  selector: 'app-row-label',
  imports: [
    MatIcon,
    NgClass
  ],
  template: `
    <div class="flex justify-center items-center gap-4 mb-2">
      @if (icon()) {
        <mat-icon [ngClass]="iconClass()">{{ icon() }}</mat-icon>
      }
      <div [ngClass]="{'grid-cols-4': layoutMode() === 'normal', 'grid-cols-2': layoutMode() === 'spaced'}" class="grid w-full items-center gap-3">
        <p class="!font-bold !text-sm whitespace-nowrap">{{ label() }}:</p>
        <div [ngClass]="{'col-span-3': layoutMode() === 'normal', 'justify-self-end': layoutMode() === 'spaced'}" class="col-span-3">
          <ng-content></ng-content>
        </div>
      </div>
    </div>


  `,
  styles: ``
})
export class RowLabelComponent {
  icon = input('')
  iconClass = input<string>('text-primary');
  label = input.required<string>()
  layoutMode = input<'tight' | 'normal' | 'spaced'>('normal');

}

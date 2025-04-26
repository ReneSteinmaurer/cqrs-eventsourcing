import { Component, input } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import {NgClass} from '@angular/common';

@Component({
  selector: 'app-status-chip',
  standalone: true,
  imports: [MatIconModule, NgClass],
  template: `
    <div [ngClass]="chipClasses()" class="px-3 py-1 rounded-full text-sm font-semibold flex items-center gap-1 w-fit">
      @if (icon()) {
        <mat-icon [color]="color()" class="text-base">{{ icon() }}</mat-icon>
      }
      {{ label() }}
    </div>
  `,
  styles: ``
})
export class StatusChipComponent {
  status = input.required<'verfuegbar' | 'verliehen' | 'katalogisiert' | 'erworben'>();

  label() {
    switch (this.status()) {
      case 'verfuegbar': return 'Verfügbar';
      case 'verliehen': return 'Verliehen';
      case 'katalogisiert': return 'Katalogisiert';
      case 'erworben': return 'Erworben';
      default: return '';
    }
  }

  color() {
    switch (this.status()) {
      case 'verfuegbar': return 'primary';
      case 'verliehen': return 'warn';
      case 'katalogisiert': return 'accent';
      case 'erworben': return undefined;
      default: return undefined;
    }
  }

  icon() {
    switch (this.status()) {
      case 'verfuegbar': return 'check_circle';
      case 'verliehen': return 'error';
      case 'katalogisiert': return 'library_books';
      case 'erworben': return 'inventory_2';
      default: return undefined;
    }
  }

  chipClasses() {
    switch (this.status()) {
      case 'verfuegbar':
        return ['bg-green-700/60', 'text-green-100', 'border', 'border-green-400/50'];
      case 'verliehen':
        return ['bg-red-700/60', 'text-red-100', 'border', 'border-red-400/50'];
      case 'katalogisiert':
        return ['bg-purple-700/60', 'text-purple-100', 'border', 'border-purple-400/50'];
      case 'erworben':
        return ['bg-gray-700/60', 'text-gray-300', 'border', 'border-gray-400/40'];
      default:
        return ['bg-white/10', 'text-white', 'border', 'border-white/20'];
    }
  }
}

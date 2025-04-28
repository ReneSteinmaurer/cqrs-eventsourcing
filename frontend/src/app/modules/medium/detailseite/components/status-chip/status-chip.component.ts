import {Component, computed, input} from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import {NgClass} from '@angular/common';
import {MediumStatus} from '../../types/medium-details';

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
  status = input.required<MediumStatus>();
  label = computed(() => {
    switch (this.status()) {
      case 'VERLIEHEN': return 'Verliehen';
      case 'VERLOREN': return 'Verloren';
      case 'KATALOGISIERT': return 'Katalogisiert';
      case 'ERWORBEN': return 'Erworben';
      default: return '';
    }
  });
  color = computed(() => {
    switch (this.status()) {
      case 'VERLOREN': return 'warn';
      case 'VERLIEHEN': return 'warn';
      case 'KATALOGISIERT': return 'accent';
      case 'ERWORBEN': return undefined;
      default: return undefined;
    }
  })
  icon = computed(() => {
    switch (this.status()) {
      case 'VERLIEHEN': return 'error';
      case 'KATALOGISIERT': return 'library_books';
      case 'ERWORBEN': return 'inventory_2';
      default: return undefined;
    }
  })
  chipClasses = computed(() => {
    switch (this.status()) {
      case 'VERLIEHEN':
        return ['!bg-[var(--mat-sys-error)]/20 ', 'text-[var(--mat-sys-error)] ', 'border', '!border-[var(--mat-sys-error)]/50'];
      case 'KATALOGISIERT':
        return ['bg-[#00dddd]/10', 'text-[#00dddd]/100', 'border', 'border-[#00dddd]/60'];
      case 'ERWORBEN':
        return ['bg-gray-700/60', 'text-gray-300', 'border', 'border-gray-400/40'];
      case 'VERLOREN':
        return ['!bg-[var(--mat-sys-error)]/20 ', 'text-[var(--mat-sys-error)] ', 'border', '!border-[var(--mat-sys-error)]/50'];
      default:
        return ['bg-white/10', 'text-white', 'border', 'border-white/20'];
    }
  })
}

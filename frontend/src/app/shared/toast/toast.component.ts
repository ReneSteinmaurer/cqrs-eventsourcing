import {Component, effect, inject} from '@angular/core';
import {NgClass} from '@angular/common';
import {ToastService} from '../services/toast.service';
import {toSignal} from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-toast',
  imports: [
    NgClass
  ],
  template: `
    <div class="fixed top-5 right-5 space-y-2 z-50">
      @for (toast of toasts(); track $index) {
        <div class="w-96 shadow-lg rounded text-white relative overflow-hidden"
             [ngClass]="bgColor(toast.type)">
          <div class="p-4">
            <strong class="block">{{ toast.title }}</strong>
            <p class="text-sm">{{ toast.message }}</p>
          </div>
          <button class="absolute hover:cursor-pointer top-1 right-2 text-white text-lg" (click)="remove(toast.id)">×</button>
          <div class="absolute bottom-0 left-0 h-1"
               [ngClass]="barColor(toast.type) + (animationClassMap.get(toast.id) ? ' shrink-animate' : '')"
               [style.--duration]="(toast.duration || 5000) + 'ms'"
               style="width: 100%;">
          </div>
        </div>
      }
    </div>
  `,
  styles: `
    .shrink-animate {
      animation: shrink var(--duration, 5000ms) linear forwards;
    }

    @keyframes shrink {
      from {
        width: 100%;
      }
      to {
        width: 0;
      }
    }
  `
})
export class ToastComponent {
  toastService = inject(ToastService);
  toasts = toSignal(this.toastService.toasts$);
  animationClassMap = new Map<string | number, boolean>();

  constructor() {
    effect(() => {
      for (const toast of this.toasts() ?? []) {
        if (!this.animationClassMap.has(toast.id)) {
          this.animationClassMap.set(toast.id, false);
          setTimeout(() => this.animationClassMap.set(toast.id, true), 10);
        }
      }
    });
  }


  remove(id: string | number) {
    this.toastService.remove(id);
  }

  bgColor(type: string) {
    return {
      success: 'bg-green-500',
      error: 'bg-red-500',
      info: 'bg-blue-500',
      warning: 'bg-yellow-500'
    }[type] || 'bg-gray-700';
  }

  barColor(type: string) {
    return {
      success: 'bg-green-300',
      error: 'bg-red-300',
      info: 'bg-blue-300',
      warning: 'bg-yellow-300'
    }[type] || 'bg-gray-300';
  }

}

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastMessage {
  id: string | number;
  type: ToastType;
  title: string;
  message: string;
  duration?: number;
}

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  private _toasts = new BehaviorSubject<ToastMessage[]>([]);
  toasts$ = this._toasts.asObservable();

  show(type: ToastType, title: string, message: string, duration = 5000) {
    const toast: ToastMessage = { id: new Date().getMilliseconds(), type, title, message, duration };
    const current = this._toasts.value;
    this._toasts.next([...current, toast]);

    setTimeout(() => {
      this._toasts.next(this._toasts.value.filter(t => t.id !== toast.id));
    }, duration);
  }

  remove(id: string | number) {
    this._toasts.next(this._toasts.value.filter(t => t.id !== id));
  }
}

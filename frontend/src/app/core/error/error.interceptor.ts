import {HttpInterceptorFn, HttpResponse} from '@angular/common/http';
import {catchError, tap, throwError} from 'rxjs';
import { ToastService } from '../../shared/services/toast.service';
import { inject } from '@angular/core';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const toastService = inject(ToastService);
  return next(req).pipe(tap((event: any) =>{
    if (event instanceof HttpResponse) {
      const body = event.body;

      if (body?.errors && Array.isArray(body.errors)) {
        for (const error of body.errors) {
          toastService.show('error', 'Error', error)
        }
      }
    }
  }),
  catchError((err: any) => {
    console.log("INTERCEPTOR")
    console.log(err)
    return throwError(err);
  }));
};

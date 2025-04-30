import {Component, effect, inject} from '@angular/core';
import {FormControl, FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatError, MatFormField, MatInput, MatLabel} from '@angular/material/input';
import {MatButton} from '@angular/material/button';
import {MatDialogRef} from '@angular/material/dialog';
import {toSignal} from '@angular/core/rxjs-interop';
import {MatAutocomplete, MatAutocompleteTrigger, MatOption} from '@angular/material/autocomplete';
import {MediumDetailService} from '../../../services/medium-detail.service';
import {debounceTime} from 'rxjs';
import {ToastService} from '../../../../../../shared/services/toast.service';

@Component({
  selector: 'app-verleihen',
  imports: [
    FormsModule,
    MatError,
    MatFormField,
    MatInput,
    MatLabel,
    ReactiveFormsModule,
    MatButton,
    MatAutocompleteTrigger,
    MatAutocomplete,
    MatOption
  ],
  template: `
    <div class="p-6 max-w-full max-h-full">
      <h3 class="text-xl font-semibold m-4">Medium verleihen</h3>

      <div class="mt-8 w-full h-full flex flex-col items-center space-y-6">
        <mat-form-field appearance="outline" class="w-1/2">
          <mat-label>Kunde</mat-label>

          <input
            type="text"
            matInput
            [formControl]="customerInputId"
            [matAutocomplete]="auto"
            required
          />

          <mat-autocomplete #auto="matAutocomplete">
            @for (option of customerOptions(); track $index) {
              <mat-option [value]="option.nutzerId">
                {{ option.vorname }} {{ option.nachname }} ({{ option.email }})
              </mat-option>
            }
          </mat-autocomplete>

          @if (customerInputId.touched && customerInputId.invalid) {
            <mat-error>Der Kunde ist erforderlich</mat-error>
          }
        </mat-form-field>
      </div>

      <div class="absolute bottom-0 right-0 mb-8 space-x-2 w-1/2 pt-4">
        <button mat-button type="button" (click)="close()">Abbrechen</button>
        <button (click)="verleihen()" mat-flat-button color="primary" [disabled]="customerInputId.invalid">Verleihen
        </button>
      </div>
    </div>
  `,
  styles: ``
})
export class VerleihenComponent {
  private dialogRef = inject(MatDialogRef<VerleihenComponent>);
  detailService = inject(MediumDetailService)
  customerOptions = this.detailService.customerOptions

  customerInputId = new FormControl<string | null>('', Validators.required);
  customerInputChanged = toSignal(this.customerInputId.valueChanges.pipe(debounceTime(200)));

  constructor() {
    effect(() => {
      const input = this.customerInputChanged();
      if (!input) {
        return
      }
      this.detailService.findNutzerByEmail(input)
    });
  }

  close() {
    this.dialogRef.close();
  }

  verleihen() {
    if (!this.customerInputId.valid) {
      return
    }
    this.detailService.verleihen(this.customerInputId.value!).subscribe(() => {
      this.dialogRef.close();
    });
  }
}

import {Component, inject} from '@angular/core';
import {LandingPageService} from '../landing-page.service';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import {FormBuilder, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatError, MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {MatButton} from '@angular/material/button';
import {KatalogisiereCommand, KatalogisiereMediumService} from './katalogisiere-medium.service';

@Component({
  selector: 'app-katalogisiere-medium-dialog',
  imports: [
    ReactiveFormsModule,
    MatFormField,
    MatInput,
    MatButton,
    MatLabel,
    MatError
  ],
  template: `
    <div class="p-6 max-w-full max-h-full">
      <h3 class="text-xl font-semibold m-4">Medium katalogisieren</h3>

      <form
        [formGroup]="form"
        (ngSubmit)="submit()"
        class="mt-8 w-full h-full flex flex-col items-center space-y-6"
      >
        <mat-form-field appearance="outline" class="w-1/2">
          <mat-label>ISBN</mat-label>
          <input matInput formControlName="ISBN" required />
          @if (form.get('ISBN')?.touched && form.get('ISBN')?.invalid) {
            <mat-error>ISBN ist erforderlich.</mat-error>
          }
        </mat-form-field>

        <mat-form-field appearance="outline" class="w-1/2">
          <mat-label>Signature</mat-label>
          <input matInput formControlName="Signature" required />
          @if (form.get('Signature')?.touched && form.get('Signature')?.invalid) {
            <mat-error>Signatur ist erforderlich.</mat-error>
          }
        </mat-form-field>

        <mat-form-field appearance="outline" class="w-1/2">
          <mat-label>Standort</mat-label>
          <input matInput formControlName="Standort" required />
          @if (form.get('Standort')?.touched && form.get('Standort')?.invalid) {
            <mat-error>Standort ist erforderlich.</mat-error>
          }
        </mat-form-field>

        <mat-form-field appearance="outline" class="w-1/2">
          <mat-label>Exemplar-Code</mat-label>
          <input matInput formControlName="ExemplarCode" required />
          @if (form.get('ExemplarCode')?.touched && form.get('ExemplarCode')?.invalid) {
            <mat-error>Exemplar-Code ist erforderlich.</mat-error>
          }
        </mat-form-field>

        <div class="absolute bottom-0 right-0 mb-8 space-x-2 w-1/2 pt-4">
          <button mat-button type="button" (click)="close()">Abbrechen</button>
          <button mat-flat-button color="primary" [disabled]="form.invalid">
            Katalogisieren
          </button>
        </div>
      </form>
    </div>
  `,
  styles: ``
})
export class KatalogisiereMediumDialogComponent {
  fb = inject(FormBuilder);
  katalogisiereService = inject(KatalogisiereMediumService)
  mediumId = inject(MAT_DIALOG_DATA).mediumId;

  private dialogRef = inject(MatDialogRef<KatalogisiereMediumDialogComponent>);

  form = this.fb.group({
    ISBN: ['', Validators.required],
    Signature: ['', Validators.required],
    Standort: ['', Validators.required],
    ExemplarCode: ['', Validators.required],
  });

  submit() {
    if (this.form.valid) {
      const command: KatalogisiereCommand = {
        MediumId: this.mediumId,
        ...(this.form.value as {
          ISBN: string;
          Signature: string;
          Standort: string;
          ExemplarCode: string;
        }),
      };
      this.katalogisiereService.katalogisiereMedium(command);
      this.dialogRef.close(command);
    }
  }


  close() {
    this.dialogRef.close();
  }
}

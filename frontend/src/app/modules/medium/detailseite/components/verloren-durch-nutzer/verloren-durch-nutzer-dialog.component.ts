import {Component, inject} from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../services/medium-detail.service';
import {ConfirmationDialogComponent} from '../../../../../shared/ui/confirmation-dialog/confirmation-dialog.component';
import {ToastService} from '../../../../../shared/services/toast.service';

@Component({
  selector: 'app-verloren-durch-nutzer-dialog',
  imports: [
    ConfirmationDialogComponent
  ],
  template: `
    <confirmation-dialog (cancel)="cancel()" (confirm)="confirm()" titel="Verloren melden">
      <p #content>
        Möchtest du dieses Medium wirklich als durch den Nutzer verloren melden?
      </p>
    </confirmation-dialog>
  `,
  styles: ``
})
export class VerlorenDurchNutzerComponent {
  private dialogRef = inject(MatDialogRef<VerlorenDurchNutzerComponent>);
  detailService = inject(MediumDetailService)

  cancel() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.verlorenDurchNutzer().subscribe(() => {
      this.dialogRef.close();
    });
  }
}

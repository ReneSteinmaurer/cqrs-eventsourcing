import {Component, inject} from '@angular/core';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../services/medium-detail.service';
import {ConfirmationDialogComponent} from '../../../../../shared/ui/confirmation-dialog/confirmation-dialog.component';
import {ToastService} from '../../../../../shared/services/toast.service';

@Component({
  selector: 'app-zuruecknehmen',
  imports: [
    ConfirmationDialogComponent
  ],
  template: `
    <confirmation-dialog (cancel)="cancel()" (confirm)="confirm()" titel="Medium zurückgegeben">
      <p #content>
        Bist du sicher, dass du das Medium als zurückgenommen markieren möchtest?
      </p>
    </confirmation-dialog>
  `,
  styles: ``
})
export class ZuruecknehmenDialogComponent {
  private dialogRef = inject(MatDialogRef<ZuruecknehmenDialogComponent>);
  detailService = inject(MediumDetailService)

  cancel() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.zurueckgeben().subscribe(() => {
      this.dialogRef.close();
    });
  }
}

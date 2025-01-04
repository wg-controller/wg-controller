export interface StoreStruct {
    ConfirmDialogShow: boolean;
    ConfirmDialogTitle: string;
    ConfirmDialogText: string;
    ConfirmDialogCallback: () => void;
    SnackBarShow: boolean;
    SnackBarText: string;
    SnackBarError: boolean;
    LoggedIn: boolean;
};
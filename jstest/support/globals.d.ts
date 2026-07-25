export type Registration = {
  type: string;
  capture: boolean;
  selector: string;
  module: string;
};

declare global {
  interface Window {
    __gsxuiRegistrations: Registration[];
  }
}

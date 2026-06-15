export interface Toast {
  id: string;
  message: string;
  type: 'error' | 'success' | 'warning' | 'info';
}

class ToastService {
  messages = $state<Toast[]>([]);

  show(message: string, type: Toast['type'] = 'info', duration = 4000) {
    const id = Math.random().toString(36).substring(2, 9);
    // Push a new toast into the messages array reactive state
    this.messages = [...this.messages, { id, message, type }];

    setTimeout(() => {
      this.dismiss(id);
    }, duration);
  }

  dismiss(id: string) {
    this.messages = this.messages.filter(m => m.id !== id);
  }

  success(message: string, duration?: number) {
    this.show(message, 'success', duration);
  }

  error(message: string, duration?: number) {
    this.show(message, 'error', duration);
  }

  warning(message: string, duration?: number) {
    this.show(message, 'warning', duration);
  }

  info(message: string, duration?: number) {
    this.show(message, 'info', duration);
  }
}

export const toastService = new ToastService();

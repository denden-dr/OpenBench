import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import CustomerDeviceDetails from './CustomerDeviceDetails.svelte';

describe('CustomerDeviceDetails Component', () => {
  it('should render customer device fields correctly', () => {
    const { getByLabelText } = render(CustomerDeviceDetails, {
      props: {
        customerName: 'John Doe',
        customerPhone: '08123456789',
        brandPhone: 'Apple',
        modelPhone: 'iPhone 15',
        serialNumber: 'IMEI123456',
        isSubmitting: false,
        isEditing: false
      }
    });

    const nameInput = getByLabelText('Customer Name') as HTMLInputElement;
    const phoneInput = getByLabelText('Phone Number') as HTMLInputElement;
    const brandInput = getByLabelText('Brand') as HTMLInputElement;
    const modelInput = getByLabelText('Model') as HTMLInputElement;
    const serialInput = getByLabelText('IMEI / Serial Number') as HTMLInputElement;

    expect(nameInput.value).toBe('John Doe');
    expect(phoneInput.value).toBe('08123456789');
    expect(brandInput.value).toBe('Apple');
    expect(modelInput.value).toBe('iPhone 15');
    expect(serialInput.value).toBe('IMEI123456');
  });

  it('should disable all inputs when isEditing is false', () => {
    const { getByLabelText } = render(CustomerDeviceDetails, {
      props: {
        customerName: 'John Doe',
        customerPhone: '08123456789',
        brandPhone: 'Apple',
        modelPhone: 'iPhone 15',
        serialNumber: 'IMEI123456',
        isSubmitting: false,
        isEditing: false
      }
    });

    expect(getByLabelText('Customer Name')).toBeDisabled();
    expect(getByLabelText('Phone Number')).toBeDisabled();
    expect(getByLabelText('Brand')).toBeDisabled();
    expect(getByLabelText('Model')).toBeDisabled();
    expect(getByLabelText('IMEI / Serial Number')).toBeDisabled();
  });

  it('should enable inputs when isEditing is true', () => {
    const { getByLabelText } = render(CustomerDeviceDetails, {
      props: {
        customerName: 'John Doe',
        customerPhone: '08123456789',
        brandPhone: 'Apple',
        modelPhone: 'iPhone 15',
        serialNumber: 'IMEI123456',
        isSubmitting: false,
        isEditing: true
      }
    });

    expect(getByLabelText('Customer Name')).not.toBeDisabled();
    expect(getByLabelText('Phone Number')).not.toBeDisabled();
    expect(getByLabelText('Brand')).not.toBeDisabled();
    expect(getByLabelText('Model')).not.toBeDisabled();
    expect(getByLabelText('IMEI / Serial Number')).not.toBeDisabled();
  });

  it('should disable inputs when isSubmitting is true even if isEditing is true', () => {
    const { getByLabelText } = render(CustomerDeviceDetails, {
      props: {
        customerName: 'John Doe',
        customerPhone: '08123456789',
        brandPhone: 'Apple',
        modelPhone: 'iPhone 15',
        serialNumber: 'IMEI123456',
        isSubmitting: true,
        isEditing: true
      }
    });

    expect(getByLabelText('Customer Name')).toBeDisabled();
    expect(getByLabelText('Phone Number')).toBeDisabled();
    expect(getByLabelText('Brand')).toBeDisabled();
    expect(getByLabelText('Model')).toBeDisabled();
    expect(getByLabelText('IMEI / Serial Number')).toBeDisabled();
  });
});

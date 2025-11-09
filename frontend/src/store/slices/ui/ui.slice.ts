import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
  duration?: number;
}

export interface UIState {
  toasts: Toast[];
  isLoading: boolean;
  sidebarOpen: boolean;
  theme: 'light' | 'dark';
}

const initialState: UIState = {
  toasts: [],
  isLoading: false,
  sidebarOpen: true,
  theme: 'light',
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    addToast: (state, action: PayloadAction<Omit<Toast, 'id'>>) => {
      const id = `${Date.now()}-${Math.random()}`;
      state.toasts.push({
        ...action.payload,
        id,
      });
    },

    removeToast: (state, action: PayloadAction<string>) => {
      state.toasts = state.toasts.filter((toast) => toast.id !== action.payload);
    },

    clearToasts: (state) => {
      state.toasts = [];
    },

    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },

    toggleSidebar: (state) => {
      state.sidebarOpen = !state.sidebarOpen;
    },

    setSidebarOpen: (state, action: PayloadAction<boolean>) => {
      state.sidebarOpen = action.payload;
    },

    setTheme: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.theme = action.payload;
    },
  },
});

export const {
  addToast,
  removeToast,
  clearToasts,
  setLoading,
  toggleSidebar,
  setSidebarOpen,
  setTheme,
} = uiSlice.actions;

export default uiSlice.reducer;

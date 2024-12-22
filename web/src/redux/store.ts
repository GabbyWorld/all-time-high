import { configureStore } from '@reduxjs/toolkit'
import { homeReducer } from './reducer'

export const store = configureStore({
  reducer: {
    homeReducer
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false
    })
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

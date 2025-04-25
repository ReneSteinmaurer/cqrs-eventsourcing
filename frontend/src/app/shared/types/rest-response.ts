export type RestResponse<T> = {
  data: T;
  errors: string[]
  message: string
}

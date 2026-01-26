import axios, { type AxiosRequestConfig } from 'axios';

export const AXIOS_INSTANCE = axios.create({
  baseURL: '/api/v1',
  withCredentials: true,
});

export const customInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const source = axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
    cancelToken: source.token,
  })
    .then(({ data }) => data)
    .catch((error) => {
      if (error.response?.status === 401 && window.location.pathname !== "/login") {
        AXIOS_INSTANCE.post("/auth/logout").then(() => {
          window.location.href = '/login';
        });
        return Promise.reject(error);
      }

      if (error.response?.data) {
        throw error.response.data;
      }

      throw error;
    });

  // @ts-ignore
  promise.cancel = () => {
    source.cancel('Query was cancelled');
  };

  return promise;
};

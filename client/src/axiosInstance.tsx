import axios from 'axios';

// Create an axios instance
const axiosInstance = axios.create({
    baseURL: 'http://localhost:8080/api/v1', // Adjust for your backend URL
    withCredentials: true, // Send credentials (cookies)
});

// Add request interceptor to attach JWT token to requests
axiosInstance.interceptors.request.use(
    (config) => {
        // Retrieve the JWT access token from cookies
        const accessToken = document.cookie.match(/(^| )access_token=([^;]+)/)?.[2]; // Or use cookies library
        if (accessToken) {
            config.headers['Authorization'] = `Bearer ${accessToken}`; // Attach access token
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Add response interceptor to handle token expiration and refresh
axiosInstance.interceptors.response.use(
    (response) => response, // If the request is successful, just return the response
    async (error) => {
        const originalRequest = error.config;

        // If we receive a 401 Unauthorized error and the original request hasn't been retried yet
        if (error.response && error.response.status === 401 && !originalRequest._retry) {
            console.log("Token Expired.");
            originalRequest._retry = true; // Mark the request as retried

            try {
                // Automatically refresh the token using the refresh token
                const refreshToken = document.cookie.match(/(^| )refresh_token=([^;]+)/)?.[2]; // Get refresh token from cookies

                if (!refreshToken) {
                    throw new Error("No refresh token available");
                }

                const refreshTokenResponse = await axiosInstance.post('/auth/refresh', { 
                    "refresh_token": refreshToken,
                    // withCredentials: true, 
                    // headers: { 
                    //     "Authorization": `Bearer ${refreshToken}` // Send the refresh token for refreshing
                    // }
                });

                const { access_token, refresh_token } = refreshTokenResponse.data;
                console.log(access_token, refresh_token)

                // Store the new tokens (in cookies)
                // document.cookie = `access_token=${access_token}; path=/;`;
                // document.cookie = `refresh_token=${refresh_token}; path=/;`;

                // Update the original request with the new access token
                originalRequest.headers['Authorization'] = `Bearer ${access_token}`;

                // Retry the original request with the new token
                return axiosInstance(originalRequest);
            } catch (refreshError) {
                console.error("Token refresh failed", refreshError);
                // window.location.href = '/login'; // Handle refresh failure (logout or redirect to login)
                return Promise.reject(refreshError);
            }
        }

        return Promise.reject(error); // Return any other errors as-is
    }
);

export default axiosInstance;

export const isAuthenticated = async (): Promise<boolean> => {
    try {
      const response = await axiosInstance.get("/users"); // Endpoint to validate token
      console.log(response)
      return response.status === 200;
    } catch (error) {
      return false; // If the request fails (e.g., 401 Unauthorized), return false
    }
  };
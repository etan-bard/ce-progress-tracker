import axios from 'axios'
import type { Takes } from '@/types'

// Define the backend response type
interface ParticipantCourseResponse {
  ParticipantID: number
  CourseID: number
  CourseName?: string
  DateFirstAccessed?: string
  DateLastAccessed: string
  CourseCompletion: number
}

// Create Axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add request interceptor for logging and potential auth tokens
api.interceptors.request.use(
  (config) => {
    console.log(`Making ${config.method?.toUpperCase()} request to ${config.url}`)
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

export const fetchParticipantCourses = async (): Promise<Takes[]> => {
  try {
    const response = await api.get<ParticipantCourseResponse[]>('/participant-courses')
    
    return response.data.map(item => ({
      participantId: item.ParticipantID,
      courseId: item.CourseID,
      courseName: item.CourseName || `Course ${item.CourseID}`,
      dateFirstAccessed: item.DateFirstAccessed || '',
      dateLastAccessed: item.DateLastAccessed,
      courseCompletion: item.CourseCompletion,
    }))
  } catch (error) {
    console.error('Failed to fetch participant courses:', error)
    throw new Error('Failed to fetch participant courses. Please try again later.')
  }
}

export default api
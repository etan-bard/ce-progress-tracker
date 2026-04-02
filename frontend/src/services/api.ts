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

// Base URL for API requests
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Helper function for logging requests
const logRequest = (method: string, url: string) => {
  console.log(`Making ${method.toUpperCase()} request to ${url}`)
}

// Helper function for error handling
const handleError = (error: unknown) => {
  console.error('API Error:', error instanceof Error ? error.message : String(error))
  throw new Error('Failed to fetch participant courses. Please try again later.')
}

export const fetchParticipantCourses = async (): Promise<Takes[]> => {
  const url = `${BASE_URL}/participant-courses`
  
  try {
    logRequest('GET', url)
    const response = await fetch(url)
    
    if (!response.ok) {
      console.error('Failed to fetch participant courses:', response.statusText)
      throw new Error('Failed to fetch participant courses. Please try again later.')
    }
    
    const data: ParticipantCourseResponse[] = await response.json()
    
    if (!data || data.length === 0) {
      console.log('No participant course data available')
      throw new Error('No participant course data available')
    }

    return data.map(item => ({
      participantId: item.ParticipantID,
      courseId: item.CourseID,
      courseName: item.CourseName || `Course ${item.CourseID}`,
      dateFirstAccessed: item.DateFirstAccessed || '',
      dateLastAccessed: item.DateLastAccessed,
      courseCompletion: item.CourseCompletion,
    }))
  } catch (error) {
    handleError(error)
    // This line is theoretically unreachable due to handleError throwing
    throw error
  }
}

export default { fetchParticipantCourses }
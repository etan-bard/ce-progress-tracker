import type {Takes} from '@/types'

// Base URL for API requests
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'


export const fetchParticipantCourses = async (): Promise<Takes[]> => {
  const url = `${BASE_URL}/participant-courses`

  const response = await fetch(url).catch(error => {
    console.error(`Network error during GET request to ${url}:`, error)
    throw new Error(`Network error during GET request to ${url}`)
  })

  if (!response.ok) {
    const errorMsg = `HTTP Error ${response.status}: ${response.statusText} for GET ${url}`
    console.error(errorMsg)
    throw new Error(errorMsg)
  }

  const data: Takes[] = await response.json().catch(error => {
    console.error(`JSON parsing error for response from ${url}:`, error)
    throw new Error(`JSON parsing error for response from ${url}`)
  })

  if (!data || data.length === 0) {
    console.log('No participant course data available')
    return []
  }

  return data.map((item) => ({
    ...item,
    courseName: `Course ${item.courseId}`,
    dateFirstAccessed: item.dateFirstAccessed || 'N/A',
    dateLastAccessed: item.dateLastAccessed || 'N/A',
  }))
}
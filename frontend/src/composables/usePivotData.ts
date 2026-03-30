import { computed, shallowRef, watch } from 'vue'
import type { Ref } from 'vue'
import type { Takes } from '@/types'

interface GridCell {
  completion: number | null
  lastAccessed: string | null
  exists: boolean
}

interface PivotData {
  participants: readonly number[]
  courses: readonly number[]
  grid: readonly (readonly GridCell[])[]
}

export const usePivotData = (linearData: Ref<Takes[]>) => {
  const pivotData = shallowRef<PivotData>({ 
    participants: Object.freeze([]), 
    courses: Object.freeze([]), 
    grid: Object.freeze([]) 
  })

  // Use computed property for automatic updates when source data changes
  const transformedData = computed(() => {
    if (!linearData.value?.length) return pivotData.value

    // Extract unique participants and courses from EXISTING data
    const participants = [...new Set(linearData.value.map(item => item.participantId))]
    const courses = [...new Set(linearData.value.map(item => item.courseId))]

    // Create grid with participantId × courseId mapping
    // Use Object.freeze to prevent unnecessary reactivity overhead
    const grid = participants.map(participantId => 
      courses.map(courseId => {
        const item = linearData.value.find(i => 
          i.participantId === participantId && i.courseId === courseId
        )
        return Object.freeze({
          completion: item?.courseCompletion || null,
          lastAccessed: item?.dateLastAccessed || null,
          exists: !!item
        })
      })
    )

    return { 
      participants: Object.freeze(participants), 
      courses: Object.freeze(courses), 
      grid: Object.freeze(grid) 
    }
  })

  // Update pivot data when transformation completes
  watch(transformedData, (newData) => {
    pivotData.value = newData
  }, { immediate: true })

  return { pivotData }
}

export const getCellValue = (pivotData: PivotData, participantId: number, courseId: number): GridCell => {
  const participantIndex = pivotData.participants.indexOf(participantId)
  const courseIndex = pivotData.courses.indexOf(courseId)
  
  if (participantIndex === -1 || courseIndex === -1 || 
      !pivotData.grid[participantIndex] || 
      pivotData.grid[participantIndex][courseIndex] === undefined) {
    return Object.freeze({ completion: null, lastAccessed: null, exists: false })
  }
  
  return pivotData.grid[participantIndex][courseIndex]
}

export type { PivotData, GridCell }
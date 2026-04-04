import { computed, type Ref } from 'vue'
import type { Takes } from '@/types'

export interface GridCell {
  completion: number | null
  lastAccessed: string | null
}

export interface PivotData {
  participants: readonly number[]
  courses: readonly number[]
  grid: readonly (readonly GridCell[])[]
}

export interface SortOption {
  key: string
  order: 'asc' | 'desc'
}

export const pivotDataSort = (a: any, b: any, isDescending: boolean) => {
  // Blank (-2) is ALWAYS last
  if (a === -2 && b === -2) return 0
  if (a === -2) return isDescending ? -1 : 1
  if (b === -2) return isDescending ? 1 : -1

  // All other values (N/A=-1 and Valid>=0)
  // Treated such that N/A < 0 < 1 ...
  return a < b ? -1 : (a > b ? 1 : 0)
}

export const usePivotData = (
  linearData: Ref<Takes[]>,
  search: Ref<string>,
  sortBy: Ref<SortOption[]>,
  sortMode: Ref<'completion' | 'date'>
) => {
  // Use computed property for automatic updates when source data changes
  const pivotData = computed<PivotData>(() => {
    if (!linearData.value?.length) {
      return {
        participants: Object.freeze([]),
        courses: Object.freeze([]),
        grid: Object.freeze([])
      }
    }

    // Extract unique participants and courses from EXISTING data
    const participants = [...new Set(linearData.value.map(item => item.participantId))].sort((a, b) => a - b)
    const courses = [...new Set(linearData.value.map(item => item.courseId))].sort((a, b) => a - b)

    // Create grid with participantId × courseId mapping
    const grid = participants.map(participantId => 
      courses.map(courseId => {
        const item = linearData.value.find(i => 
          i.participantId === participantId && i.courseId === courseId
        )
        return Object.freeze({
          completion: item?.courseCompletion ?? null,
          lastAccessed: item?.dateLastAccessed ?? null
        })
      })
    )

    return { 
      participants: Object.freeze(participants), 
      courses: Object.freeze(courses), 
      grid: Object.freeze(grid) 
    }
  })

  // Computed property for the final table items, filtered and sorted
  const transformedItems = computed(() => {
    const { participants, courses } = pivotData.value
    
    // 1. Create participant objects with all course data
    let items = participants.map(participantId => {
      const participantObj: Record<string, any> = {
        participantId: participantId
      }
      
      courses.forEach(courseId => {
        const cellValue = getCellValue(pivotData.value, participantId, courseId)
        participantObj[courseId.toString()] = cellValue
        
        // Handle completion % for sorting
        participantObj[`${courseId}_completion`] = cellValue?.completion ?? -2 // -2 represents null
        
        // Handle last accessed date for sorting
        if (cellValue?.lastAccessed) {
          const timestamp = new Date(cellValue.lastAccessed).getTime()
          // Dates around 1969/1970 (unix epoch 0) should be treated as N/A (-1)
          // 31536000000 is 1 year in milliseconds (Jan 1 1971)
          participantObj[`${courseId}_date`] = timestamp < 31536000000 ? -1 : timestamp
        } else if (cellValue?.completion !== null) {
          // completion exists but date does not => N/A
          participantObj[`${courseId}_date`] = -1
        } else {
          // No record at all => null
          participantObj[`${courseId}_date`] = -2
        }
      })
      
      return participantObj
    })

    // 2. Filter by search
    if (search.value) {
      const searchLower = search.value.toLowerCase()
      items = items.filter(item => 
        item.participantId.toString().includes(searchLower)
      )
    }

    // 3. Sort items
    if (sortBy.value.length > 0) {
      const { key, order } = sortBy.value[0] || { key: '', order: 'asc' };
      const isDescending = order === 'desc'

      items.sort((a, b) => {
        if (key === 'participantId') {
          const valA = a.participantId
          const valB = b.participantId
          const cmp = valA < valB ? -1 : (valA > valB ? 1 : 0)
          return isDescending ? -cmp : cmp
        }
        
        // For course columns, we let v-data-table handle the sorting via the custom 'sort' function in dynamicHeaders.
        // However, we still perform a basic sort here to keep the items state consistent if needed,
        // but we MUST use the same logic as the v-data-table custom sort.
        const courseId = key
        const keySuffix = sortMode.value === 'completion' ? 'completion' : 'date'
        const valA = a[`${courseId}_${keySuffix}`]
        const valB = b[`${courseId}_${keySuffix}`]

        if (valA === undefined || valB === undefined) return 0

        const cmp = pivotDataSort(valA, valB, isDescending)
        return isDescending ? -cmp : cmp
      })
    }

    return items
  })

  return { pivotData, transformedItems }
}

export const getCellValue = (pivotData: PivotData, participantId: number, courseId: number): GridCell => {
  const participantIndex = pivotData.participants.indexOf(participantId)
  const courseIndex = pivotData.courses.indexOf(courseId)
  
  if (participantIndex === -1 || courseIndex === -1 || 
      !pivotData.grid[participantIndex] || 
      pivotData.grid[participantIndex][courseIndex] === undefined) {
    return Object.freeze({ completion: null, lastAccessed: null })
  }
  
  return pivotData.grid[participantIndex][courseIndex]
}
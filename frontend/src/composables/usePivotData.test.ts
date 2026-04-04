import { describe, it, expect } from 'vitest'
import { usePivotData, type SortOption, pivotDataSort } from './usePivotData'
import { ref } from 'vue'
import type { Takes } from '@/types'

describe('pivotDataSort', () => {
  it('should always keep Blank (-2) at the end in ASC', () => {
    expect(pivotDataSort(-2, 10, false)).toBe(1)
    expect(pivotDataSort(10, -2, false)).toBe(-1)
    expect(pivotDataSort(-2, -2, false)).toBe(0)
  })

  it('should always keep Blank (-2) at the end in DESC', () => {
    expect(pivotDataSort(-2, 10, true)).toBe(-1)
    expect(pivotDataSort(10, -2, true)).toBe(1)
    expect(pivotDataSort(-2, -2, true)).toBe(0)
  })

  it('should treat N/A (-1) as smaller than any valid number', () => {
    expect(pivotDataSort(-1, 0, false)).toBeLessThan(0)
    expect(pivotDataSort(0, -1, false)).toBeGreaterThan(0)
    expect(pivotDataSort(-1, 100, false)).toBeLessThan(0)
  })
})

describe('usePivotData', () => {
  const mockTakes: Takes[] = [
    { participantId: 1, courseId: 101, courseCompletion: 100, dateLastAccessed: '2023-01-01', dateFirstAccessed: '2023-01-01' },
    { participantId: 1, courseId: 102, courseCompletion: 50, dateLastAccessed: '2023-01-01', dateFirstAccessed: '2023-01-01' },
    { participantId: 2, courseId: 101, courseCompletion: 0, dateLastAccessed: '2023-01-01', dateFirstAccessed: '2023-01-01' },
    { participantId: 2, courseId: 102, courseCompletion: 0, dateLastAccessed: '1970-01-01', dateFirstAccessed: '1970-01-01' }, // N/A (epoch)
    { participantId: 3, courseId: 101, courseCompletion: 75, dateLastAccessed: '2023-01-01', dateFirstAccessed: '2023-01-01' },
    { participantId: 3, courseId: 102, courseCompletion: 10, dateLastAccessed: '', dateFirstAccessed: '' }, // N/A (no date)
    // Participant 4 has no record for Course 102 (null/blanks)
    { participantId: 4, courseId: 101, courseCompletion: 10, dateLastAccessed: '2023-01-01', dateFirstAccessed: '2023-01-01' },
  ]

  it('should sort Course 102 by date in ASC order correctly', () => {
    const linearData = ref(mockTakes)
    const search = ref('')
    const sortBy = ref<SortOption[]>([{ key: '102', order: 'asc' }])
    const sortMode = ref<'completion' | 'date'>('date')

    const { transformedItems } = usePivotData(linearData, search, sortBy, sortMode)
    
    const results = transformedItems.value.map(item => item.participantId)

    expect(results[results.length - 1]).toBe(4) // blank, should be last
    expect(results[2]).toBe(1) // valid date, should be after N/As
    expect([2, 3]).toContain(results[0])
    expect([2, 3]).toContain(results[1])
  })

  it('should sort Course 102 by date in DESC order correctly', () => {
    const linearData = ref(mockTakes)
    const search = ref('')
    const sortBy = ref<SortOption[]>([{ key: '102', order: 'desc' }])
    const sortMode = ref<'completion' | 'date'>('date')

    const { transformedItems } = usePivotData(linearData, search, sortBy, sortMode)
    
    const results = transformedItems.value.map(item => item.participantId)
    // DESC: Valid dates first (highest first), then N/A (-1), then Blanks (-2)
    // P1 (timestamp), [P2, P3] (-1), P4 (-2)
    
    expect(results[1]).toBe(1)
    expect([2, 3]).toContain(results[2])
    expect(results[3]).toBe(4)
  })

  it('should sort Course 102 by completion in ASC order correctly', () => {
    const linearData = ref(mockTakes)
    const search = ref('')
    const sortBy = ref<SortOption[]>([{ key: '102', order: 'asc' }])
    const sortMode = ref<'completion' | 'date'>('completion')

    const { transformedItems } = usePivotData(linearData, search, sortBy, sortMode)
    
    const results = transformedItems.value.map(item => item.participantId)

    // ASC: 0 (P2), 10 (P3), 50 (P1), Blank (P4)
    expect(results[0]).toBe(2)
    expect(results[1]).toBe(3)
    expect(results[2]).toBe(1)
    expect(results[3]).toBe(4)
  })

  it('should sort Course 102 by completion in DESC order correctly', () => {
    const linearData = ref(mockTakes)
    const search = ref('')
    const sortBy = ref<SortOption[]>([{ key: '102', order: 'desc' }])
    const sortMode = ref<'completion' | 'date'>('completion')

    const { transformedItems } = usePivotData(linearData, search, sortBy, sortMode)
    
    const results = transformedItems.value.map(item => item.participantId)
    // DESC: 50 (P1), 10 (P3), 0 (P2), Blank (P4)
    expect(results[0]).toBe(1)
    expect(results[1]).toBe(3)
    expect(results[2]).toBe(2)
    expect(results[3]).toBe(4)
  })
})

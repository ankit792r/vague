/** Keep wildmenu on one line; pan without showing a scrollbar. */
export function syncWildmenuScroll(
  menu: HTMLDivElement,
  items: readonly (HTMLElement | null)[],
  index: number,
): void {
  const selected = items[index]
  if (!selected) {
    menu.scrollLeft = 0
    return
  }

  const padding = 4
  const menuWidth = menu.clientWidth
  if (menuWidth <= 0) {
    return
  }

  let scrollLeft = menu.scrollLeft
  const viewLeft = scrollLeft
  const viewRight = scrollLeft + menuWidth

  const selLeft = selected.offsetLeft
  const selRight = selLeft + selected.offsetWidth

  if (selLeft < viewLeft) {
    scrollLeft = Math.max(0, selLeft - padding)
  } else if (selRight > viewRight) {
    scrollLeft = selRight - menuWidth + padding
  }

  menu.scrollLeft = scrollLeft

  const visibleIndices: number[] = []
  for (let i = 0; i < items.length; i++) {
    const el = items[i]
    if (!el) {
      continue
    }
    const left = el.offsetLeft
    const right = left + el.offsetWidth
    if (right > scrollLeft && left < scrollLeft + menuWidth) {
      visibleIndices.push(i)
    }
  }

  const pos = visibleIndices.indexOf(index)
  if (pos === -1) {
    return
  }

  const lastVisibleSlot = Math.max(0, visibleIndices.length - 3)
  if (pos < lastVisibleSlot) {
    return
  }

  const last = items[items.length - 1]
  if (!last) {
    return
  }

  const lastRight = last.offsetLeft + last.offsetWidth
  const currentRight = scrollLeft + menuWidth
  if (lastRight > currentRight) {
    menu.scrollLeft = Math.max(0, lastRight - menuWidth + padding)
  }
}

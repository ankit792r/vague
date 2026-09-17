/** Measure how many monospace cells fit in the editor viewport. */
export function measureEditor(el: HTMLElement) {
  const style = getComputedStyle(el)
  const fontSize = parseFloat(style.fontSize) || 16
  const lineHeight = parseFloat(style.lineHeight)
  const effectiveLineHeight =
    Number.isFinite(lineHeight) && lineHeight > 0 ? lineHeight : fontSize * 1.4

  const padX = parseFloat(style.paddingLeft) + parseFloat(style.paddingRight)
  const padY = parseFloat(style.paddingTop) + parseFloat(style.paddingBottom)
  const contentWidth = Math.max(0, el.clientWidth - padX)
  const contentHeight = Math.max(0, el.clientHeight - padY)

  const charWidth = measureCharWidth(el)
  if (charWidth <= 0) {
    return { rows: 1, cols: 1 }
  }

  return {
    rows: Math.max(1, Math.floor(contentHeight / effectiveLineHeight)),
    cols: Math.max(1, Math.floor(contentWidth / charWidth)),
  }
}

function measureCharWidth(el: HTMLElement): number {
  const style = getComputedStyle(el)
  const canvas = document.createElement("canvas")
  const ctx = canvas.getContext("2d")
  if (!ctx) {
    return fontSizeFallback(style)
  }

  ctx.font = `${style.fontStyle} ${style.fontVariant} ${style.fontWeight} ${style.fontSize} ${style.fontFamily}`

  // Average width over digits — stable for monospace faces like Iosevka.
  const sample = "0123456789"
  return ctx.measureText(sample).width / sample.length
}

function fontSizeFallback(style: CSSStyleDeclaration): number {
  return parseFloat(style.fontSize) || 16
}

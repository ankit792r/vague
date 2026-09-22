import type { EditorTabState } from "../../types/editor"

type TabLineProps = {
  tabs: EditorTabState[]
}

export function TabLine({ tabs }: TabLineProps) {
  if (tabs.length <= 1) {
    return null
  }

  return (
    <div class="tab-line" aria-label="tab line">
      {tabs.map((tab, i) => (
        <span
          key={`${tab.label}-${i}`}
          class={`tab-label${tab.active ? " tab-label-active" : ""}`}
        >
          {tab.label}
        </span>
      ))}
    </div>
  )
}

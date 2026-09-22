import { describe, expect, it } from "vitest"
import { encodeKey } from "./keys"

describe("encodeKey", () => {
  it("maps control letters", () => {
    const e = {
      key: "k",
      code: "KeyK",
      ctrlKey: true,
      altKey: false,
    } as KeyboardEvent
    expect(encodeKey(e)).toBe("<C-k>")
  })

  it("uses code for letter keys when key is not length 1", () => {
    const e = {
      key: "Dead",
      code: "KeyA",
      ctrlKey: false,
      altKey: false,
    } as KeyboardEvent
    expect(encodeKey(e)).toBe("a")
  })
})

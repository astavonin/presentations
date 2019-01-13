class: center, middle

# Implicit and explicit dependencies
## ??? Subtitle of document ???

23 Jan 2018,
Alexander Stavonin

---

# Dependencies types
- Explicit dependency
```go
r, tn, p := collectFromSelectors(t)
res = append(r, sel.Sel.Name)
pos = p
typeName = tn
```

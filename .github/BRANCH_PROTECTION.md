# Branch Protection for `main`

To protect the `main` branch on GitHub, go to:

**Settings → Branches → Add branch protection rule**

### Recommended Settings:

| Setting | Value |
|---------|-------|
| Branch name pattern | `main` |
| Require a pull request before merging | ✅ |
| Require approvals | 1 |
| Require status checks to pass before merging | ✅ |
| Required status checks | `Build & Test`, `Docker Build` |
| Require branches to be up to date before merging | ✅ |
| Require conversation resolution before merging | ✅ |
| Do not allow bypassing the above settings | ✅ (optional for solo projects) |

### Quick Setup via GitHub CLI:

```bash
gh api repos/Golangcodes/golangcodes/branches/main/protection \
  -X PUT \
  -H "Accept: application/vnd.github+json" \
  -f "required_status_checks[strict]=true" \
  -f "required_status_checks[contexts][]=Build & Test" \
  -f "required_status_checks[contexts][]=Docker Build" \
  -f "required_pull_request_reviews[required_approving_review_count]=1" \
  -f "enforce_admins=false" \
  -f "restrictions=null"
```

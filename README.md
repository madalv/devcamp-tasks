### Task 1

All SQL statements can be found in `statements.sql`.

To create the database, simply run `make run`, then `make migup`, which will run the migrations on the MariaDB database.

`db/migrations` currently holds the first migration that defines the initial schema.

The code that seeds the database is in `util/seeder.go`.

### Task 2

Set up a simple router (`api/handler.go`) and a very simple cache (`repository/local_cache.go`). 
There is a benchmark for the `getCampaignsForSource` handler in `api/handler_test.go`.

Run `make bench` to run the test.

Result for the benchmark with the cache:
```
  657496              1649 ns/op
PASS
ok      adt/api 2.129s
```
Results for the benchmark without the cache:
```
    1858            647353 ns/op
PASS
ok      adt/api 1.275s
```

### Task 3

Added whitelist/blacklist fields for campaign (currently just a list of strings). Addded a `filterCampaigns` function that 
checks if the query domain is contained (or a subdomain of a domain contained) in a campaign's whitelist or blacklist.

### Task 4



#### Benchmark (original)

##### memprofile

Below are some functions detailing the memory usage when filtering campaigns. The surprise here, I 
suppose, is that `regexp.MatchString` is using dramatically more memory than anything else.
```
(pprof) list filterCampaigns
Total: 2.67GB
ROUTINE ======================== adt/api.filterCampaigns in /home/madalv/adtelligent/api/handlers.go
         0     2.61GB (flat, cum) 97.98% of Total
         .          .     91:func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
         .          .     92:   if domain == "" {
         .          .     93:           return camps
         .          .     94:   } else {
         .          .     95:           domain = strings.ToLower(domain)
         .          .     96:   }
         .          .     97:
         .          .     98:   for _, c := range camps {
         .     2.61GB     99:           contained := domainInList(domain, c.DomainList)
         .          .    100:
         .          .    101:           if (contained && c.ListType == model.BLACKLIST) ||
         .          .    102:                   (!contained && c.ListType == model.WHITELIST) {
         .          .    103:                   //slog.Debug("Campaign skipped", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
         .          .    104:                   continue
(pprof) 
(pprof) list domainInList
Total: 2.67GB
ROUTINE ======================== adt/api.domainInList in /home/madalv/adtelligent/api/handlers.go
   20.50MB     2.61GB (flat, cum) 97.98% of Total
         .          .    118:func domainInList(queryDomain string, list []string) bool {
         .          .    119:   for _, domain := range list {
         .          .    120:           domain = strings.ToLower(domain)
   20.50MB    20.50MB    121:           regexPattern := `(^|\.)(` + domain + `)($)`
         .          .    122:           // check if the queryDomain is a subdomain of current list item
         .     2.59GB    123:           match, err := regexp.MatchString(regexPattern, queryDomain)
         .          .    124:           if err != nil {
         .          .    125:                   return false
         .          .    126:           }
         .          .    127:           if match {
         .          .    128:                   return true

(pprof) top
Showing nodes accounting for 2.61GB, 97.87% of 2.67GB total
Dropped 17 nodes (cum <= 0.01GB)
Showing top 10 nodes out of 33
      flat  flat%   sum%        cum   cum%
    1.55GB 58.10% 58.10%     1.55GB 58.10%  regexp/syntax.(*compiler).inst (inline)
    0.66GB 24.73% 82.82%     0.66GB 24.73%  regexp/syntax.(*parser).newRegexp (inline)
    0.09GB  3.35% 86.17%     0.09GB  3.35%  regexp/syntax.(*parser).maybeConcat
    0.07GB  2.71% 88.88%     2.59GB 96.99%  regexp.compile
    0.06GB  2.27% 91.15%     0.91GB 34.24%  regexp/syntax.parse
    0.06GB  2.16% 93.31%     0.14GB  5.34%  regexp/syntax.(*parser).push
    0.05GB  1.74% 95.05%     0.15GB  5.71%  regexp/syntax.(*parser).collapse
    0.03GB  1.15% 96.20%     0.03GB  1.15%  regexp/syntax.(*Regexp).CapNames
    0.02GB  0.88% 97.08%     0.04GB  1.63%  github.com/brianvoe/gofakeit/v7.domainName
    0.02GB  0.79% 97.87%     0.04GB  1.52%  regexp/syntax.(*compiler).init (inline)
```

##### cpufprofile

Again, regexes seem to take a lot of CPU time, as well as functions used in garbage collection.

```
(pprof) top
Showing nodes accounting for 2230ms, 39.68% of 5620ms total
Dropped 100 nodes (cum <= 28.10ms)
Showing top 10 nodes out of 137
      flat  flat%   sum%        cum   cum%
     420ms  7.47%  7.47%      650ms 11.57%  runtime.findObject
     260ms  4.63% 12.10%      770ms 13.70%  regexp/syntax.(*compiler).inst
     260ms  4.63% 16.73%     1180ms 21.00%  runtime.mallocgc
     240ms  4.27% 21.00%     1510ms 26.87%  runtime.scanobject
     180ms  3.20% 24.20%      180ms  3.20%  runtime.memclrNoHeapPointers
     180ms  3.20% 27.40%      180ms  3.20%  runtime.nextFreeFast (inline)
     180ms  3.20% 30.60%      180ms  3.20%  runtime.spanOf (inline)
     180ms  3.20% 33.81%      260ms  4.63%  runtime.typePointers.next
     170ms  3.02% 36.83%      170ms  3.02%  runtime.memmove
     160ms  2.85% 39.68%      310ms  5.52%  regexp/syntax.(*parser).maybeConcat
(pprof) list filterCampaigns
Total: 5.62s
ROUTINE ======================== adt/api.filterCampaigns in /home/madalv/adtelligent/api/handlers.go
         0      3.67s (flat, cum) 65.30% of Total
         .          .     91:func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
         .          .     92:   if domain == "" {
         .          .     93:           return camps
         .          .     94:   } else {
         .          .     95:           domain = strings.ToLower(domain)
         .          .     96:   }
         .          .     97:
         .          .     98:   for _, c := range camps {
         .      3.67s     99:           contained := domainInList(domain, c.DomainList)
         .          .    100:
         .          .    101:           if (contained && c.ListType == model.BLACKLIST) ||
         .          .    102:                   (!contained && c.ListType == model.WHITELIST) {
         .          .    103:                   //slog.Debug("Campaign skipped", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
         .          .    104:                   continue
(pprof) list domainInList
Total: 5.62s
ROUTINE ======================== adt/api.domainInList in /home/madalv/adtelligent/api/handlers.go
      10ms      3.67s (flat, cum) 65.30% of Total
         .          .    118:func domainInList(queryDomain string, list []string) bool {
      10ms       10ms    119:   for _, domain := range list {
         .       70ms    120:           domain = strings.ToLower(domain)
         .       50ms    121:           regexPattern := `(^|\.)(` + domain + `)($)`
         .          .    122:           // check if the queryDomain is a subdomain of current list item
         .      3.54s    123:           match, err := regexp.MatchString(regexPattern, queryDomain)
         .          .    124:           if err != nil {
         .          .    125:                   return false
         .          .    126:           }
         .          .    127:           if match {
         .          .    128:                   return true

```

#### Benchmark (using domain map instead of list + removed regexes)

##### memprofile

There is a significant decrease in memory usage once both regexes are removed and domains are contained in maps instead of slices.

```
(pprof) top
Showing nodes accounting for 1859.57MB, 100% of 1860.08MB total
Dropped 5 nodes (cum <= 9.30MB)
Showing top 10 nodes out of 16
      flat  flat%   sum%        cum   cum%
  596.37MB 32.06% 32.06%  1115.38MB 59.96%  adt/api.filterCampaigns
  329.68MB 17.72% 49.79%   744.18MB 40.01%  adt/api.generateCampaigns
  306.01MB 16.45% 66.24%   306.01MB 16.45%  strings.(*Builder).grow
  293.01MB 15.75% 81.99%   293.01MB 15.75%  strings.genSplit
     229MB 12.31% 94.30%   414.51MB 22.28%  github.com/brianvoe/gofakeit/v7.domainName
  105.50MB  5.67%   100%   105.50MB  5.67%  fmt.Sprintf
         0     0%   100%  1859.57MB   100%  adt/api.BenchmarkFilterCampaigns
         0     0%   100%   519.01MB 27.90%  adt/api.domainInList
         0     0%   100%   414.51MB 22.28%  github.com/brianvoe/gofakeit/v7.DomainName (inline)
         0     0%   100%   306.01MB 16.45%  strings.(*Builder).Grow
(pprof) list filterCampaigns
Total: 1.82GB
ROUTINE ======================== adt/api.filterCampaigns in /home/madalv/adtelligent/api/handlers.go
  596.37MB     1.09GB (flat, cum) 59.96% of Total
         .          .     90:func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
         .          .     91:   if domain == "" {
         .          .     92:           return camps
         .          .     93:   } else {
         .          .     94:           domain = strings.ToLower(domain)
         .          .     95:   }
         .          .     96:
         .          .     97:   for _, c := range camps {
         .   519.01MB     98:           contained := domainInList(domain, c.DomainList)
         .          .     99:
         .          .    100:           if (contained && c.ListType == model.BLACKLIST) ||
         .          .    101:                   (!contained && c.ListType == model.WHITELIST) {
         .          .    102:                   //slog.Debug("Campaign skipped", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
         .          .    103:                   continue
         .          .    104:           }
         .          .    105:           //slog.Debug("Campaign good", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
  596.37MB   596.37MB    106:           filtered = append(filtered, c)
         .          .    107:   }
         .          .    108:   return

```

##### cpuprofile

A reduction in CPU usage is also noticed, especially in the function `domainInList`, since it 
doesn't have to iterate over slices or match regex patterns. CPU seems to be used mostly by runtime functions for accessing maps and allocating memory.

```
(pprof) top
Showing nodes accounting for 2840ms, 49.74% of 5710ms total
Dropped 95 nodes (cum <= 28.55ms)
Showing top 10 nodes out of 113
      flat  flat%   sum%        cum   cum%
     550ms  9.63%  9.63%      840ms 14.71%  runtime.mapaccess2_faststr
     350ms  6.13% 15.76%      470ms  8.23%  runtime.mapaccess1_faststr
     340ms  5.95% 21.72%     1200ms 21.02%  runtime.mallocgc
     280ms  4.90% 26.62%      280ms  4.90%  aeshashbody
     250ms  4.38% 31.00%      410ms  7.18%  strings.ToLower
     230ms  4.03% 35.03%      230ms  4.03%  runtime.memclrNoHeapPointers
     220ms  3.85% 38.88%      220ms  3.85%  runtime.memmove
     220ms  3.85% 42.73%      890ms 15.59%  runtime.scanobject
     200ms  3.50% 46.23%      270ms  4.73%  runtime.findObject
     200ms  3.50% 49.74%      200ms  3.50%  runtime.nextFreeFast (inline)
(pprof) list filterCampaigns
Total: 5.71s
ROUTINE ======================== adt/api.filterCampaigns in /home/madalv/adtelligent/api/handlers.go
      30ms      1.62s (flat, cum) 28.37% of Total
         .          .     90:func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
         .          .     91:   if domain == "" {
         .          .     92:           return camps
         .          .     93:   } else {
         .          .     94:           domain = strings.ToLower(domain)
         .          .     95:   }
         .          .     96:
      10ms       10ms     97:   for _, c := range camps {
         .      1.38s     98:           contained := domainInList(domain, c.DomainList)
         .          .     99:
         .          .    100:           if (contained && c.ListType == model.BLACKLIST) ||
         .          .    101:                   (!contained && c.ListType == model.WHITELIST) {
         .          .    102:                   //slog.Debug("Campaign skipped", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
         .          .    103:                   continue
         .          .    104:           }
         .          .    105:           //slog.Debug("Campaign good", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
      20ms      230ms    106:           filtered = append(filtered, c)
         .          .    107:   }
         .          .    108:   return
         .          .    109:}
         .          .    110:
         .          .    111:/*
(pprof) list domainInList
Total: 5.71s
ROUTINE ======================== adt/api.domainInList in /home/madalv/adtelligent/api/handlers.go
      10ms      1.38s (flat, cum) 24.17% of Total
         .          .    117:func domainInList(queryDomain string, dMap map[string]struct{}) bool {
         .      440ms    118:   parts := strings.Split(queryDomain, ".")
      10ms       10ms    119:   for i := 0; i < len(parts)-1; i++ {
         .      410ms    120:           currDomain := strings.Join(parts[i:], ".")
         .      520ms    121:           _, ok := dMap[currDomain]
         .          .    122:           if ok {
         .          .    123:                   return true
         .          .    124:           }
         .          .    125:   }
         .          .    126:

```

import type { DefineComponent, SlotsType } from 'vue'
type IslandComponent<T> = DefineComponent<{}, {refresh: () => Promise<void>}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, SlotsType<{ fallback: { error: unknown } }>> & T

type HydrationStrategies = {
  hydrateOnVisible?: IntersectionObserverInit | true
  hydrateOnIdle?: number | true
  hydrateOnInteraction?: keyof HTMLElementEventMap | Array<keyof HTMLElementEventMap> | true
  hydrateOnMediaQuery?: string
  hydrateAfter?: number
  hydrateWhen?: boolean
  hydrateNever?: true
}
type LazyComponent<T> = DefineComponent<HydrationStrategies, {}, {}, {}, {}, {}, {}, { hydrated: () => void }> & T

interface _GlobalComponents {
  Badge: typeof import("../../app/components/ui/badge/index")['Badge']
  Button: typeof import("../../app/components/ui/button/index")['Button']
  Card: typeof import("../../app/components/ui/card/index")['Card']
  CardAction: typeof import("../../app/components/ui/card/index")['CardAction']
  CardContent: typeof import("../../app/components/ui/card/index")['CardContent']
  CardDescription: typeof import("../../app/components/ui/card/index")['CardDescription']
  CardFooter: typeof import("../../app/components/ui/card/index")['CardFooter']
  CardHeader: typeof import("../../app/components/ui/card/index")['CardHeader']
  CardTitle: typeof import("../../app/components/ui/card/index")['CardTitle']
  Input: typeof import("../../app/components/ui/input/index")['Input']
  Select: typeof import("../../app/components/ui/select/index")['Select']
  SelectContent: typeof import("../../app/components/ui/select/index")['SelectContent']
  SelectGroup: typeof import("../../app/components/ui/select/index")['SelectGroup']
  SelectItem: typeof import("../../app/components/ui/select/index")['SelectItem']
  SelectItemText: typeof import("../../app/components/ui/select/index")['SelectItemText']
  SelectLabel: typeof import("../../app/components/ui/select/index")['SelectLabel']
  SelectScrollDownButton: typeof import("../../app/components/ui/select/index")['SelectScrollDownButton']
  SelectScrollUpButton: typeof import("../../app/components/ui/select/index")['SelectScrollUpButton']
  SelectSeparator: typeof import("../../app/components/ui/select/index")['SelectSeparator']
  SelectTrigger: typeof import("../../app/components/ui/select/index")['SelectTrigger']
  SelectValue: typeof import("../../app/components/ui/select/index")['SelectValue']
  Separator: typeof import("../../app/components/ui/separator/index")['Separator']
  Table: typeof import("../../app/components/ui/table/index")['Table']
  TableBody: typeof import("../../app/components/ui/table/index")['TableBody']
  TableCaption: typeof import("../../app/components/ui/table/index")['TableCaption']
  TableCell: typeof import("../../app/components/ui/table/index")['TableCell']
  TableEmpty: typeof import("../../app/components/ui/table/index")['TableEmpty']
  TableFooter: typeof import("../../app/components/ui/table/index")['TableFooter']
  TableHead: typeof import("../../app/components/ui/table/index")['TableHead']
  TableHeader: typeof import("../../app/components/ui/table/index")['TableHeader']
  TableRow: typeof import("../../app/components/ui/table/index")['TableRow']
  TableUtils: typeof import("../../app/components/ui/table/utils")['default']
  NuxtWelcome: typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']
  NuxtLayout: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
  NuxtErrorBoundary: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
  ClientOnly: typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']
  DevOnly: typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']
  ServerPlaceholder: typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']
  NuxtLink: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']
  NuxtLoadingIndicator: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
  NuxtTime: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
  NuxtRouteAnnouncer: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
  NuxtImg: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
  NuxtPicture: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
  NuxtPage: typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']
  NoScript: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']
  Link: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']
  Base: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']
  Title: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']
  Meta: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']
  Style: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']
  Head: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']
  Html: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']
  Body: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']
  NuxtIsland: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']
  LazyBadge: LazyComponent<typeof import("../../app/components/ui/badge/index")['Badge']>
  LazyButton: LazyComponent<typeof import("../../app/components/ui/button/index")['Button']>
  LazyCard: LazyComponent<typeof import("../../app/components/ui/card/index")['Card']>
  LazyCardAction: LazyComponent<typeof import("../../app/components/ui/card/index")['CardAction']>
  LazyCardContent: LazyComponent<typeof import("../../app/components/ui/card/index")['CardContent']>
  LazyCardDescription: LazyComponent<typeof import("../../app/components/ui/card/index")['CardDescription']>
  LazyCardFooter: LazyComponent<typeof import("../../app/components/ui/card/index")['CardFooter']>
  LazyCardHeader: LazyComponent<typeof import("../../app/components/ui/card/index")['CardHeader']>
  LazyCardTitle: LazyComponent<typeof import("../../app/components/ui/card/index")['CardTitle']>
  LazyInput: LazyComponent<typeof import("../../app/components/ui/input/index")['Input']>
  LazySelect: LazyComponent<typeof import("../../app/components/ui/select/index")['Select']>
  LazySelectContent: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectContent']>
  LazySelectGroup: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectGroup']>
  LazySelectItem: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectItem']>
  LazySelectItemText: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectItemText']>
  LazySelectLabel: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectLabel']>
  LazySelectScrollDownButton: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectScrollDownButton']>
  LazySelectScrollUpButton: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectScrollUpButton']>
  LazySelectSeparator: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectSeparator']>
  LazySelectTrigger: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectTrigger']>
  LazySelectValue: LazyComponent<typeof import("../../app/components/ui/select/index")['SelectValue']>
  LazySeparator: LazyComponent<typeof import("../../app/components/ui/separator/index")['Separator']>
  LazyTable: LazyComponent<typeof import("../../app/components/ui/table/index")['Table']>
  LazyTableBody: LazyComponent<typeof import("../../app/components/ui/table/index")['TableBody']>
  LazyTableCaption: LazyComponent<typeof import("../../app/components/ui/table/index")['TableCaption']>
  LazyTableCell: LazyComponent<typeof import("../../app/components/ui/table/index")['TableCell']>
  LazyTableEmpty: LazyComponent<typeof import("../../app/components/ui/table/index")['TableEmpty']>
  LazyTableFooter: LazyComponent<typeof import("../../app/components/ui/table/index")['TableFooter']>
  LazyTableHead: LazyComponent<typeof import("../../app/components/ui/table/index")['TableHead']>
  LazyTableHeader: LazyComponent<typeof import("../../app/components/ui/table/index")['TableHeader']>
  LazyTableRow: LazyComponent<typeof import("../../app/components/ui/table/index")['TableRow']>
  LazyTableUtils: LazyComponent<typeof import("../../app/components/ui/table/utils")['default']>
  LazyNuxtWelcome: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
  LazyNuxtLayout: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
  LazyNuxtErrorBoundary: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
  LazyClientOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']>
  LazyDevOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']>
  LazyServerPlaceholder: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
  LazyNuxtLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
  LazyNuxtLoadingIndicator: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
  LazyNuxtTime: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
  LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
  LazyNuxtImg: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
  LazyNuxtPicture: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
  LazyNuxtPage: LazyComponent<typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']>
  LazyNoScript: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
  LazyLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']>
  LazyBase: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']>
  LazyTitle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']>
  LazyMeta: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']>
  LazyStyle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']>
  LazyHead: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']>
  LazyHtml: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']>
  LazyBody: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']>
  LazyNuxtIsland: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']>
}

declare module 'vue' {
  export interface GlobalComponents extends _GlobalComponents { }
}

export {}
